package writers

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"metricsd/config"
	"metricsd/types"
	"github.com/kpumuk/gorrd"
)

type rrdUpdateTask struct {
	writer         Writer
	firstSampleSet *types.SampleSet
	firstDataItem  dataItem
	f              func([]string) []string
	wg             *sync.WaitGroup
}

var (
	// Channel with tasks for RRD update threads
	rrdUpdateTasks chan *rrdUpdateTask
	// Indicating whether RRD update threads were created
	rrdUpdateThreadsPrepared bool = false
)

func Rollup(writer Writer, set *types.SampleSet) {
	prepareRrdUpdateThreads()
	wg := &sync.WaitGroup{}

	if data := writer.rollupData(set); data != nil {
		updateRrd(writer, set, data, wg, func(args []string) []string {
			return append(args, data.rrdString())
		})
	}

	wg.Wait()
}

func BatchRollup(writer Writer, sets []*types.SampleSet) {
	data := make([]dataItem, 0, 10)

	var from int
	var prevSource, prevName string

	prepareRrdUpdateThreads()
	wg := &sync.WaitGroup{}

	for cur, set := range sets {
		// config.Logger.Debug("... source=%s, name=%s, prevSource=%s, prevName=%s", set.Source, set.Name, prevSource, prevName)
		if cur == 0 {
			prevSource = set.Source
			prevName = set.Name
		}

		// Next item in the sequence of samples
		pushed := false
		if prevSource == set.Source && prevName == set.Name {
			if item := writer.rollupData(set); item != nil {
				data = append(data, item)
			}
			pushed = true
		}

		// Reached a new sequence or the end of samples list
		if prevSource != set.Source || prevName != set.Name || cur == len(sets)-1 {
			batchRollup(writer, sets[from], data, wg)

			from = cur
			prevSource = set.Source
			prevName = set.Name
			data = make([]dataItem, 0, 10)
		}

		// A new sequence beginning
		if !pushed {
			if item := writer.rollupData(set); item != nil {
				data = append(data, item)

				// The last item in the samples list
				if cur == len(sets)-1 {
					batchRollup(writer, sets[from], data, wg)
				}
			}
		}
	}

	wg.Wait()
}

func batchRollup(writer Writer, firstSampleSet *types.SampleSet, data []dataItem, wg *sync.WaitGroup) {
	// Nothing to save
	if len(data) == 0 {
		return
	}

	// Update RRD database
	updateRrd(writer, firstSampleSet, data[0], wg, func(args []string) []string {
		// Serialize all data items to the arguments array
		for _, elem := range data {
			args = append(args, elem.rrdString())
		}
		return args
	})
}

func prepareRrdUpdateThreads() {
	if rrdUpdateThreadsPrepared {
		return
	}

	rrdUpdateTasks = make(chan *rrdUpdateTask, config.RrdUpdateThreads)
	for i := 1; i <= config.RrdUpdateThreads; i++ {
		go func(idx int) {
			config.Logger.Debug("Started RRD update thread #%d", idx)
			args := make([]string, 0, 10)
			runtime.LockOSThread()
			for {
				task := <-rrdUpdateTasks
				args = task.f(args[:0])
				doUpdateRrd(task.writer, task.firstSampleSet, task.firstDataItem, args)
				task.wg.Done()
			}
		}(i)
	}
	rrdUpdateThreadsPrepared = true
}

func updateRrd(writer Writer, firstSampleSet *types.SampleSet, firstDataItem dataItem, wg *sync.WaitGroup, f func([]string) []string) {
	wg.Add(1)
	rrdUpdateTasks <- &rrdUpdateTask{writer: writer, firstSampleSet: firstSampleSet, firstDataItem: firstDataItem, f: f, wg: wg}
}

func doUpdateRrd(writer Writer, firstSampleSet *types.SampleSet, firstDataItem dataItem, args []string) {
	file := getRrdFile(writer, firstSampleSet)
	if _, err := os.Stat(file); err != nil {
		err := rrd.Create(file, int64(config.SliceInterval), firstSampleSet.Time-int64(config.SliceInterval), firstDataItem.rrdInfo())
		if err != nil {
			config.Logger.Debug("Error occurred: %s", err)
			return
		}
	}
	// config.Logger.Debug("... file=%s", file)
	err := rrd.Update(file, firstDataItem.rrdTemplate(), args)
	if err != nil {
		config.Logger.Debug("Error occurred: %s", err)
	}
}

func getRrdFile(writer Writer, set *types.SampleSet) string {
	dir := fmt.Sprintf("%s/%s", config.DataDir, set.Source)
	os.MkdirAll(dir, 0755)
	// This is temporary solution while we migrate from $ grouping to .
	file := fmt.Sprintf("%s-%s", strings.Replace(set.Name, "$", ".", -1), writer.Name())
	path := fmt.Sprintf("%s/%s.rrd", dir, file)
	migrateDollarGroupsToDots(dir, file, path)
	return path
}

func migrateDollarGroupsToDots(dir, file, path string) {
	if _, err := os.Stat(path); err != nil {
		oldFile := strings.Replace(file, ".", "$", 1)
		oldPath := fmt.Sprintf("%s/%s.rrd", dir, oldFile)
		// config.Logger.Info("Probing %s", oldPath)
		if _, err := os.Stat(oldPath); err == nil {
			config.Logger.Info("Old file exists, renaming %s to %s", oldPath, path)
			os.Rename(oldPath, path)
			return
		}

		oldFile = strings.Replace(oldFile, ".", "_", -1)
		oldPath = fmt.Sprintf("%s/%s.rrd", dir, oldFile)
		// config.Logger.Info("Probing %s", oldPath)
		if _, err := os.Stat(oldPath); err == nil {
			config.Logger.Info("Old file exists, renaming %s to %s", oldPath, path)
			os.Rename(oldPath, path)
			return
		}

		oldFile = strings.Replace(oldFile, "$", ".", -1)
		oldPath = fmt.Sprintf("%s/%s.rrd", dir, oldFile)
		// config.Logger.Info("Probing %s", oldPath)
		if _, err := os.Stat(oldPath); err == nil {
			config.Logger.Info("Old file exists, renaming %s to %s", oldPath, path)
			os.Rename(oldPath, path)
			return
		}

		oldFile = strings.Replace(oldFile, ".", "_", -1)
		oldPath = fmt.Sprintf("%s/%s.rrd", dir, oldFile)
		// config.Logger.Info("Probing %s", oldPath)
		if _, err := os.Stat(oldPath); err == nil {
			config.Logger.Info("Old file exists, renaming %s to %s", oldPath, path)
			os.Rename(oldPath, path)
			return
		}
	}
}

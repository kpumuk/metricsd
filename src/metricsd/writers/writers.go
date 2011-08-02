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
	writer Writer
	set *types.SampleSet
	data dataItem
	f func() []string
	wg *sync.WaitGroup
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
		updateRrd(writer, set, data, wg, func() []string {
			return []string{data.rrdString()}
		})
	}

	wg.Wait()
}

func BatchRollup(writer Writer, sets []*types.SampleSet) {
	data := make([]dataItem, 0, len(sets))
	args := make([]string, 0, len(sets))

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
			batchRollup(writer, from, sets, data, &args, wg)

			from = cur
			prevSource = set.Source
			prevName = set.Name
			data = data[0:0]
		}

		// A new sequence beginning
		if !pushed {
			if item := writer.rollupData(set); item != nil {
				data = append(data, item)

				// The last item in the samples list
				if cur == len(sets)-1 {
					batchRollup(writer, from, sets, data, &args, wg)
				}
			}
		}
	}

	wg.Wait()
}

func batchRollup(writer Writer, from int, sets []*types.SampleSet, data []dataItem, buf *[]string, wg *sync.WaitGroup) {
	// config.Logger.Debug("Starting batchRollup [time=%v, writer=%T, from=%d, name=%s, source=%s]", time.Nanoseconds(), writer, from, sets[from].Name, sets[from].Source)
	// Nothing to save
	if len(data) == 0 {
		return
	}

	// Retrieve the first data item (used to get RRD-related information)
	firstItem := data[0]
	// Retrieve the first sample set (used to generate RRD file name)
	firstSet := sets[from]

	// Data collector function
	f := func() (args []string) {
		// Serialize all data items to buffer
		args = (*buf)[:len(data)]
		for i, elem := range data {
			args[i] = elem.rrdString()
		}
		// config.Logger.Debug("... args=%s", args)
		return
	}

	// Update RRD database
	updateRrd(writer, firstSet, firstItem, wg, f)
}

func prepareRrdUpdateThreads() {
	if rrdUpdateThreadsPrepared {
		return
	}

	rrdUpdateTasks = make(chan *rrdUpdateTask, config.RrdUpdateThreads)
	for i := 1; i <= config.RrdUpdateThreads; i++ {
		go func(idx int) {
			config.Logger.Debug("Started RRD update thread #%d", idx)
			runtime.LockOSThread()
			for {
				task := <-rrdUpdateTasks
				doUpdateRrd(task.writer, task.set, task.data, task.f)
				task.wg.Done()
			}
		}(i)
	}
	rrdUpdateThreadsPrepared = true
}

func updateRrd(writer Writer, set *types.SampleSet, data dataItem, wg *sync.WaitGroup, f func() []string) {
	wg.Add(1)
	rrdUpdateTasks <- &rrdUpdateTask{writer: writer, set: set, data: data, f: f, wg: wg}
}

func doUpdateRrd(writer Writer, set *types.SampleSet, data dataItem, f func() []string) {
	file := getRrdFile(writer, set)
	if _, err := os.Stat(file); err != nil {
		err := rrd.Create(file, int64(config.SliceInterval), set.Time-int64(config.SliceInterval), data.rrdInfo())
		if err != nil {
			config.Logger.Debug("Error occurred: %s", err)
			return
		}
	}
	// config.Logger.Debug("... file=%s", file)
	err := rrd.Update(file, data.rrdTemplate(), f())
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

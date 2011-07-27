package writers

import (
	"fmt"
	"os"
	"metricsd/config"
	"metricsd/types"
	"github.com/kpumuk/gorrd"
)


type Writer interface {
	Name() string
	Rollup(set *types.SampleSet)
	BatchRollup(sets types.SampleSetsList)
	// Private methods
	rollupData(set *types.SampleSet) dataItem
}

type dataItem interface {
	rrdInfo() []string
	rrdTemplate() string
	rrdString() string
	String() string
}

func Rollup(writer Writer, set *types.SampleSet) {
	if data := writer.rollupData(set); data != nil {
		updateRrd(writer, set, data, func() []string {
			return []string{data.rrdString()}
		})
	}
}

func BatchRollup(writer Writer, sets types.SampleSetsList) {
	data := make([]dataItem, 0, len(sets))
	args := make([]string, 0, len(sets))

	var from int
	var prevSource, prevName string

	for cur, set := range sets {
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
		if prevSource != set.Source || prevName != set.Name || cur == sets.Len()-1 {
			batchRollup(writer, from, sets, data, &args)

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
				if cur == sets.Len()-1 {
					batchRollup(writer, from, sets, data, &args)
				}
			}
		}
	}
}

func batchRollup(writer Writer, from int, sets types.SampleSetsList, data []dataItem, buf *[]string) {
	// Nothing to save
	if len(data) == 0 {
		return
	}

	// Retrieve the first data item (used to get RRD-related information)
	firstItem := data[0].(dataItem)
	// Retrieve the first sample set (used to generate RRD file name)
	firstSet := sets[from]
	// Update RRD database
	updateRrd(writer, firstSet, firstItem, func() (args []string) {
		// Serialize all data items to buffer
		args = (*buf)[0:len(data)]
		for i, elem := range data {
			args[i] = elem.(dataItem).rrdString()
		}
		return
	})
}

func updateRrd(writer Writer, set *types.SampleSet, data dataItem, f func() []string) {
	file := getRrdFile(writer, set)
	if _, err := os.Stat(file); err != nil {
		err := rrd.Create(file, int64(config.SliceInterval), set.Time-int64(config.SliceInterval), data.rrdInfo())
		if err != nil {
			config.Logger.Debug("Error occurred: %s", err)
			return
		}
	}
	err := rrd.Update(file, data.rrdTemplate(), f())
	if err != nil {
		config.Logger.Debug("Error occurred: %s", err)
	}
}

func getRrdFile(writer Writer, set *types.SampleSet) string {
	dir := fmt.Sprintf("%s/%s", config.DataDir, set.Source)
	os.MkdirAll(dir, 0755)
	return fmt.Sprintf("%s/%s-%s.rrd", dir, set.Name, writer.Name())
}

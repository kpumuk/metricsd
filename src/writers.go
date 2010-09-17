package writers

import (
    "container/vector"
    "fmt"
    "os"
    "rrd"
    "./config"
    "./types"
)


type Writer interface {
    Name() string
    Rollup(set *types.SampleSet)
    BatchRollup(sets *vector.Vector)
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
            return []string { data.rrdString() }
        })
    }
}

func BatchRollup(writer Writer, sets *vector.Vector) {
    data := new(vector.Vector).Resize(0, sets.Len())
    args := make([]string, 0, sets.Len())

    var from int
    var prevSource, prevName string

    for cur, elem := range *sets {
        set := elem.(*types.SampleSet)
        if cur == 0 { prevSource = set.Source; prevName = set.Name }

        // Next item in the sequence of samples
        pushed := false
        if prevSource == set.Source && prevName == set.Name {
            if item := writer.rollupData(set); item != nil {
                data.Push(&item)
            }
            pushed = true
        }

        // Reached a new sequence or the end of samples list
        if prevSource != set.Source || prevName != set.Name || cur == sets.Len() - 1 {
            batchRollup(writer, from, sets, data, &args)

            from       = cur
            prevSource = set.Source
            prevName   = set.Name
            data.Resize(0, 0)
        }

        // A new sequence beginning
        if !pushed {
            if item := writer.rollupData(set); item != nil {
                data.Push(&item)

                // The last item in the samples list
                if cur == sets.Len() - 1 {
                    batchRollup(writer, from, sets, data, &args)
                }
            }
        }
    }
}

func batchRollup(writer Writer, from int, sets *vector.Vector, data *vector.Vector, buf *[]string) {
    // Nothing to save
    if data.Len() == 0 { return }

    // Retrieve the first data item (used to get RRD-related information)
    firstItem := data.At(0).(*dataItem)
    // Retrieve the first sample set (used to generate RRD file name)
    firstSet  := sets.At(from).(*types.SampleSet)
    // Update RRD database
    updateRrd(writer, firstSet, *firstItem, func() (args []string) {
        // Serialize all data items to buffer
        args = (*buf)[0:data.Len()]
        for i, elem := range *data {
            args[i] = elem.(*dataItem).rrdString()
        }
        return
    })
}

func updateRrd(writer Writer, set *types.SampleSet, data dataItem, f func() []string) {
    file := getRrdFile(writer, set)
    if _, err := os.Stat(file); err != nil {
        err := rrd.Create(file, int64(config.Global.SliceInterval), set.Time - int64(config.Global.SliceInterval), data.rrdInfo())
        if err != nil {
            config.Global.Logger.Debug("Error occurred: %s", err)
            return
        }
    }
    err := rrd.Update(file, data.rrdTemplate(), f())
    if err != nil {
        config.Global.Logger.Debug("Error occurred: %s", err)
    }
}

func getRrdFile(writer Writer, set *types.SampleSet) string {
    dir := fmt.Sprintf("%s/%s", config.Global.DataDir, set.Source)
    os.MkdirAll(dir, 0755)
    return fmt.Sprintf("%s/%s-%s.rrd", dir, set.Name, writer.Name())
}

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

    var from, cur int
    var prevKey string
    var firstItem *dataItem
    var pushed bool
    var elem interface{}

    for cur, elem = range *sets {
        set := elem.(*types.SampleSet)
        if cur == 0 { prevKey = set.Key }

        pushed = false
        if prevKey == set.Key {
            if item := writer.rollupData(set); item != nil {
                data.Push(&item)
                if firstItem == nil { firstItem = &item }
            }
            pushed = true
        }

        if prevKey != set.Key || cur == sets.Len() - 1 {
            if firstItem != nil {
                firstSet := sets.At(from).(*types.SampleSet)
                updateRrd(writer, firstSet, *firstItem, func() []string {
                    ar := args[0:data.Len()]
                    for j, el := range *data {
                        ar[j] = el.(*dataItem).rrdString()
                    }
                    return ar
                })
            }

            firstItem = nil
            from      = cur
            prevKey   = set.Key
            data.Resize(0, 0)
        }

        if !pushed {
            if item := writer.rollupData(set); item != nil {
                data.Push(&item)
                if firstItem == nil { firstItem = &item }
            }
        }
    }
}

func updateRrd(writer Writer, set *types.SampleSet, data dataItem, f func() []string) {
    file := getRrdFile(writer, set)
    if _, err := os.Stat(file); err != nil {
        err := rrd.Create(file, 10, set.Time - 10, data.rrdInfo())
        if err != nil {
            config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
            return
        }
    }
    err := rrd.Update(file, data.rrdTemplate(), f())
    if err != nil {
        config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
    }
}

func getRrdFile(writer Writer, set *types.SampleSet) string {
    dir := fmt.Sprintf("%s/%s", config.GlobalConfig.DataDir, set.Source)
    os.MkdirAll(dir, 0755)
    return fmt.Sprintf("%s/%s-%s.rrd", dir, set.Name, writer.Name())
}

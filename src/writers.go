package writers

import (
    "fmt"
    "os"
    "rrd"
    "./config"
    "./types"
)


type Writer interface {
    Name() string
    Rollup(set *types.SampleSet)
    // Private methods
    rollupData(set *types.SampleSet) dataItem
}

type dataItem interface {
    rrdInfo() []string
    rrdTemplate() string
    rrdString(time int64) string
}

func Rollup(writer Writer, set *types.SampleSet) {
    if data := writer.rollupData(set); data != nil {
        updateRrd(writer, set, data)
    }
}

func updateRrd(writer Writer, set *types.SampleSet, data dataItem) {
    file := getRrdFile(writer, set)
    if _, err := os.Stat(file); err != nil {
        err := rrd.Create(file, 10, set.Time - 10, data.rrdInfo())
        if err != nil {
            config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
            return
        }
    }
    err := rrd.Update(file, data.rrdTemplate(), []string { data.rrdString(set.Time) })
    if err != nil {
        config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
    }
}

func getRrdFile(writer Writer, set *types.SampleSet) string {
    dir := fmt.Sprintf("%s/%s", config.GlobalConfig.DataDir, set.Source)
    os.MkdirAll(dir, 0755)
    return fmt.Sprintf("%s/%s-%s.rrd", dir, set.Name, writer.Name())
}

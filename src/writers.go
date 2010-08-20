package writers

import (
    "fmt"
    "os"
    "rrd"
    "sort"
    "./config"
    "./types"
)


type Writer interface {
    Rollup(set *types.SampleSet)
}

func getRrdFile(t string, set *types.SampleSet) string {
    dir := fmt.Sprintf("%s/%s", config.GlobalConfig.DataDir, set.Source)
    os.MkdirAll(dir, 0755)
    return fmt.Sprintf("%s/%s-%s.rrd", dir, set.Name, t)
}

/******************************************************************************/

type Quartiles struct {
}

type QuartilesItem struct {
    lo, q1, q2, q3, hi, total int
}

func (quartiles *Quartiles) Rollup(set *types.SampleSet) {
    if set.Values.Len() < 2 { return }
    sort.Sort(set.Values)
    number := set.Values.Len()
    lo := set.Values.At(0)
    hi := set.Values.At(number - 1)
    lo_c := number / 2
    hi_c := number - lo_c
    data := &QuartilesItem {}
    if lo_c > 0 && hi_c > 0 {
        lo_samples := set.Values.Slice(0, lo_c)
        hi_samples := set.Values.Slice(lo_c, lo_c + hi_c)
        lo_sum := 0
        hi_sum := 0
        lo_samples.Do(func(elem int) { lo_sum += elem })
        hi_samples.Do(func(elem int) { hi_sum += elem })
        q1 := lo_sum / lo_c
        q2 := (lo_sum + hi_sum) / (lo_c + hi_c)
        q3 := hi_sum / hi_c

        data.lo = lo
        data.q1 = q1
        data.q2 = q2
        data.q3 = q3
        data.hi = hi
        data.total = number

        quartiles.save(set, data)
    }
}

func (self *Quartiles) save(set *types.SampleSet, data *QuartilesItem) {
    file := getRrdFile("quartiles", set)
    if _, err := os.Stat(file); err != nil {
        err := rrd.Create(file, 10, set.Time - 10, []string {
            "DS:q1:GAUGE:600:0:U",
            "DS:q2:GAUGE:600:0:U",
            "DS:q3:GAUGE:600:0:U",
            "DS:hi:GAUGE:600:0:U",
            "DS:lo:GAUGE:600:0:U",
            "DS:total:GAUGE:600:0:U",
            "RRA:AVERAGE:0.5:1:25920",      // 72 hours at 1 sample per 10 secs
            "RRA:AVERAGE:0.5:60:4320",      // 1 month at 1 sample per 10 mins
            "RRA:AVERAGE:0.5:2880:5475",    // 5 years at 1 sample per 8 hours
            "RRA:MIN:0.5:1:25920",          // 72 hours at 1 sample per 10 secs
            "RRA:MIN:0.5:60:4320",          // 1 month at 1 sample per 10 mins
            "RRA:MIN:0.5:2880:5475",        // 5 years at 1 sample per 8 hours
            "RRA:MAX:0.5:1:25920",          // 72 hours at 1 sample per 10 secs
            "RRA:MAX:0.5:60:4320",          // 1 month at 1 sample per 10 mins
            "RRA:MAX:0.5:2880:5475",        // 5 years at 1 sample per 8 hours
        })
        if err != nil {
            config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
            return
        }
    }
    err := rrd.Update(file, "q1:q2:q3:lo:hi:total", []string {
        fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d", set.Time, data.q1, data.q2, data.q3, data.lo, data.hi, data.total),
    })
    if err != nil {
        config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
    }
}

/******************************************************************************/

type YesOrNo struct {
}

type YesOrNoItem struct {
    ok   uint64
    fail uint64
}

func (self *YesOrNo) Rollup(set *types.SampleSet) {
    data := &YesOrNoItem {}
    set.Values.Do(func(elem int) {
        if elem >= 0 {
            data.ok++
        } else {
            data.fail++
        }
    })
    self.save(set, data)
}

func (self *YesOrNo) save(set *types.SampleSet, data *YesOrNoItem) {
    file := getRrdFile("yesno", set)
    if _, err := os.Stat(file); err != nil {
        err := rrd.Create(file, 10, set.Time - 10, []string {
            "DS:ok:GAUGE:600:0:U",
            "DS:fail:GAUGE:600:0:U",
            "RRA:AVERAGE:0.5:1:25920",      // 72 hours at 1 sample per 10 secs
            "RRA:AVERAGE:0.5:60:4320",      // 1 month at 1 sample per 10 mins
            "RRA:AVERAGE:0.5:2880:5475",    // 5 years at 1 sample per 8 hours
        })

        if err != nil {
            config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
            return
        }
    }
    err := rrd.Update(file, "ok:fail", []string {
        fmt.Sprintf("%d:%d:%d", set.Time, data.ok, data.fail),
    })
    if err != nil {
        config.GlobalConfig.Logger.Debug("Error occurred: %s", err)
    }
}

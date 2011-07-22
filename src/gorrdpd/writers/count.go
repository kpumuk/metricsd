package writers

import (
    "fmt"
    "gorrdpd/types"
)

type Count struct{}

type CountItem struct {
    time int64
    ok   uint64
    fail uint64
}

func (*Count) Name() string {
    return "count"
}

func (self *Count) Rollup(set *types.SampleSet) {
    Rollup(self, set)
}

func (self *Count) BatchRollup(sets types.SampleSetsList) {
    BatchRollup(self, sets)
}

func (self *Count) rollupData(set *types.SampleSet) (data dataItem) {
    var ok, fail uint64
    for _, elem := range set.Values {
        if elem >= 0 {
            ok++
        } else {
            fail++
        }
    }
    data = &CountItem{time: set.Time, ok: ok, fail: fail}
    return
}

func (self *CountItem) String() string {
    return fmt.Sprintf("CountItem[time=%d, ok=%d, fail=%d]", self.time, self.ok, self.fail)
}

func (*CountItem) rrdInfo() []string {
    return []string{
        "DS:ok:ABSOLUTE:600:0:U",
        "DS:fail:ABSOLUTE:600:0:U",
        "RRA:AVERAGE:0.5:1:25920",   // 72 hours at 1 sample per 10 secs
        "RRA:AVERAGE:0.5:60:4320",   // 1 month at 1 sample per 10 mins
        "RRA:AVERAGE:0.5:2880:5475", // 5 years at 1 sample per 8 hours
    }
}

func (*CountItem) rrdTemplate() string {
    return "ok:fail"
}

func (self *CountItem) rrdString() string {
    return fmt.Sprintf("%d:%d:%d", self.time, self.ok, self.fail)
}

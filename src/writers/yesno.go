package writers

import (
    "container/vector"
    "fmt"
    "./types"
)

type YesOrNo struct {
}

type YesOrNoItem struct {
    time int64
    ok   uint64
    fail uint64
}

func (*YesOrNo) Name() string {
    return "yesno"
}

func (self *YesOrNo) Rollup(set *types.SampleSet) {
    Rollup(self, set)
}

func (self *YesOrNo) BatchRollup(sets *vector.Vector) {
    BatchRollup(self, sets)
}

func (self *YesOrNo) rollupData(set *types.SampleSet) (data dataItem) {
    var ok, fail uint64
    set.Values.Do(func(elem int) {
        if elem >= 0 {
            ok++
        } else {
            fail++
        }
    })
    data = &YesOrNoItem { time: set.Time, ok: ok, fail: fail }
    return
}

func (self *YesOrNoItem) String() string {
    return fmt.Sprintf("YesOrNoItem[time=%d, ok=%d, fail=%d]", self.time, self.ok, self.fail)
}

func (*YesOrNoItem) rrdInfo() []string {
    return []string {
        "DS:ok:GAUGE:600:0:U",
        "DS:fail:GAUGE:600:0:U",
        "RRA:AVERAGE:0.5:1:25920",      // 72 hours at 1 sample per 10 secs
        "RRA:AVERAGE:0.5:60:4320",      // 1 month at 1 sample per 10 mins
        "RRA:AVERAGE:0.5:2880:5475",    // 5 years at 1 sample per 8 hours
    }
}

func (*YesOrNoItem) rrdTemplate() string {
    return "ok:fail"
}

func (self *YesOrNoItem) rrdString() string {
    return fmt.Sprintf("%d:%d:%d", self.time, self.ok, self.fail)
}

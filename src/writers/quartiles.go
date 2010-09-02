package writers

import (
    "container/vector"
    "fmt"
    "sort"
    "./types"
)

type Quartiles struct {
}

type QuartilesItem struct {
    time int64
    lo, q1, q2, q3, hi, total int
}

func (self *Quartiles) Name() string {
    return "quartiles"
}

func (self *Quartiles) Rollup(set *types.SampleSet) {
    Rollup(self, set)
}

func (self *Quartiles) BatchRollup(sets *vector.Vector) {
    BatchRollup(self, sets)
}

func (self *Quartiles) rollupData(set *types.SampleSet) (data dataItem) {
    if set.Values.Len() < 2 { return }
    sort.Sort(set.Values)
    number := set.Values.Len()
    lo := set.Values.At(0)
    hi := set.Values.At(number - 1)
    lo_c := number / 2
    hi_c := number - lo_c
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

        data = &QuartilesItem {
            time: set.Time,
            lo: lo,
            q1: q1,
            q2: q2,
            q3: q3,
            hi: hi,
            total: number,
        }
    }
    return
}

// func (self *Quartiles) BatchRollup(sampleSets *vector.Vector) {
//     var from int = 0
//     var prevKey string = ""
//     for i := 0; i < sampleSets.Len(); i++ {
//         elem := sampleSets.At(i).(*types.SampleSet)
//         if i == sampleSets.Len() - 1 || (prevKey != elem.Key && prevKey != "") {
//             data = make([]*QuartilesItem, 0, i - from + 1)
//             for j := from; j <= i; j++ {
//                 data[j - from] = rollupData()
//             }
//             from = i
//         }
//     }
//
// }

func (self *QuartilesItem) String() string {
    return fmt.Sprintf(
        "QuartilesItem[time=%d, lo=%d, q1=%d, q2=%d, q3=%d, hi=%d, total=%d]",
        self.time,
        self.lo,
        self.q1,
        self.q2,
        self.q3,
        self.hi,
        self.total,
    )
}

func (*QuartilesItem) rrdInfo() []string {
    return []string {
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
    }
}

func (*QuartilesItem) rrdTemplate() string {
    return "q1:q2:q3:lo:hi:total"
}

func (self *QuartilesItem) rrdString() string {
    return fmt.Sprintf(
        "%d:%d:%d:%d:%d:%d:%d",
        self.time,
        self.q1,
        self.q2,
        self.q3,
        self.lo,
        self.hi,
        self.total,
    )
}

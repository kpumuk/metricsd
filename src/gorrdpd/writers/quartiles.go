package writers

import (
    "fmt"
    "sort"
    "gorrdpd/types"
)

type Quartiles struct{}

type QuartilesItem struct {
    time                      int64
    lo, q1, q2, q3, hi, total int
}

func (self *Quartiles) Name() string {
    return "quartiles"
}

func (self *Quartiles) Rollup(set *types.SampleSet) {
    Rollup(self, set)
}

func (self *Quartiles) BatchRollup(sets types.SampleSetsList) {
    BatchRollup(self, sets)
}

func (self *Quartiles) rollupData(set *types.SampleSet) (data dataItem) {
    if len(set.Values) < 2 {
        return
    }
    sort.Sort(set.Values)
    number := len(set.Values)
    lo := set.Values[0]
    hi := set.Values[number - 1]
    lo_c := number / 2
    hi_c := number - lo_c
    if lo_c > 0 && hi_c > 0 {
        lo_samples := set.Values[0:lo_c]
        hi_samples := set.Values[lo_c:lo_c+hi_c]
        lo_sum := 0
        hi_sum := 0
        for _, elem := range lo_samples {
	    	lo_sum += elem
		}
        for _, elem := range hi_samples {
			hi_sum += elem
		}
        q1 := lo_sum / lo_c
        q2 := (lo_sum + hi_sum) / (lo_c + hi_c)
        q3 := hi_sum / hi_c

        data = &QuartilesItem{
            time:  set.Time,
            lo:    lo,
            q1:    q1,
            q2:    q2,
            q3:    q3,
            hi:    hi,
            total: number,
        }
    }
    return
}

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
    return []string{
        "DS:q1:GAUGE:600:0:U",
        "DS:q2:GAUGE:600:0:U",
        "DS:q3:GAUGE:600:0:U",
        "DS:hi:GAUGE:600:0:U",
        "DS:lo:GAUGE:600:0:U",
        "DS:total:ABSOLUTE:600:0:U",
        "RRA:AVERAGE:0.5:1:25920",   // 72 hours at 1 sample per 10 secs
        "RRA:AVERAGE:0.5:60:4320",   // 1 month at 1 sample per 10 mins
        "RRA:AVERAGE:0.5:2880:5475", // 5 years at 1 sample per 8 hours
        "RRA:MIN:0.5:1:25920",       // 72 hours at 1 sample per 10 secs
        "RRA:MIN:0.5:60:4320",       // 1 month at 1 sample per 10 mins
        "RRA:MIN:0.5:2880:5475",     // 5 years at 1 sample per 8 hours
        "RRA:MAX:0.5:1:25920",       // 72 hours at 1 sample per 10 secs
        "RRA:MAX:0.5:60:4320",       // 1 month at 1 sample per 10 mins
        "RRA:MAX:0.5:2880:5475",     // 5 years at 1 sample per 8 hours
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

package writers

import (
	"fmt"
	"sort"
	"gorrdpd/types"
)

// Quartiles writer is used to calculate quartiles, min and max values, and
// total number of values in a sample set.
type Quartiles struct{}

type QuartilesItem struct {
	// Timestamp of the sample set.
	time int64
	// Minimum value in the sample set.
	lo int64
	// Q1 (25%)
	q1 int64
	// Q2 (50%)
	q2 int64
	// Q3 (75%)
	q3 int64
	// Maximum value in the sample set.
	hi int64
	// Number of values used to generate statistics.
	total int64
}

// Name returns the name of the writer.
func (self *Quartiles) Name() string {
	return "quartiles"
}

// Rollup performs summarization on the given sample set and writes
// results to RRD file.
func (self *Quartiles) Rollup(set *types.SampleSet) {
	Rollup(self, set)
}

// BatchRollup performs summarization on the given list of sample sets and
// writes results to RRD files.
func (self *Quartiles) BatchRollup(sets types.SampleSetsList) {
	BatchRollup(self, sets)
}

// rollupData performs summarization on the given sample set and returns
// QuartilesItem with statistics.
func (self *Quartiles) rollupData(set *types.SampleSet) (data dataItem) {
	if len(set.Values) < 2 {
		return
	}
	sort.Sort(set.Values)
	number := int64(len(set.Values))
	lo := int64(set.Values[0])
	hi := int64(set.Values[number - 1])
	lo_c := number / 2
	hi_c := number - lo_c
	if lo_c > 0 && hi_c > 0 {
		var lo_sum int64 = 0
		var hi_sum int64 = 0
		for _, elem := range set.Values[0:lo_c] {
			lo_sum += int64(elem)
		}
		for _, elem := range set.Values[lo_c:lo_c+hi_c] {
			hi_sum += int64(elem)
		}
		q1 := lo_sum / lo_c
		q2 := (lo_sum + hi_sum) / (lo_c + hi_c)
		q3 := hi_sum / hi_c

		data = &QuartilesItem{
			time:	set.Time,
			lo:		lo,
			q1:		q1,
			q2:		q2,
			q3:		q3,
			hi:		hi,
			total:	number,
		}
	}
	return
}

// String returns string representation of the given QuartilesItem.
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

// rrdInfo returns the list of parameters used to create RRD file.
func (*QuartilesItem) rrdInfo() []string {
	return []string{
		"DS:q1:GAUGE:600:0:U",
		"DS:q2:GAUGE:600:0:U",
		"DS:q3:GAUGE:600:0:U",
		"DS:hi:GAUGE:600:0:U",
		"DS:lo:GAUGE:600:0:U",
		"DS:total:ABSOLUTE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",		// 72 hours at 1 sample per 10 secs
		"RRA:AVERAGE:0.5:60:4320",		// 1 month at 1 sample per 10 mins
		"RRA:AVERAGE:0.5:2880:5475",	// 5 years at 1 sample per 8 hours
		"RRA:MIN:0.5:1:25920",			// 72 hours at 1 sample per 10 secs
		"RRA:MIN:0.5:60:4320",			// 1 month at 1 sample per 10 mins
		"RRA:MIN:0.5:2880:5475",		// 5 years at 1 sample per 8 hours
		"RRA:MAX:0.5:1:25920",			// 72 hours at 1 sample per 10 secs
		"RRA:MAX:0.5:60:4320",			// 1 month at 1 sample per 10 mins
		"RRA:MAX:0.5:2880:5475",		// 5 years at 1 sample per 8 hours
	}
}

// rrdTemplate returns template for RRDTool used to update data.
func (*QuartilesItem) rrdTemplate() string {
	return "q1:q2:q3:lo:hi:total"
}

// rrdString returns a string matching template format with the data to
// update RRD files.
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

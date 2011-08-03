package writers

import (
	"fmt"
	"math"
	"sort"
	"metricsd/types"
)

// Quartiles writer is used to calculate quartiles, min and max values, and
// total number of values in a sample set.
type Quartiles struct {
	*BaseWriter
}

type quartilesItem struct {
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

// rollupData performs summarization on the given sample set and returns
// quartilesItem with statistics.
func (self *Quartiles) rollupData(set *types.SampleSet) (data dataItem) {
	if len(set.Values) == 0 {
		return
	}
	sort.Ints(set.Values)
	number := int64(len(set.Values))
	lo := int64(set.Values[0])
	hi := int64(set.Values[number-1])

	q1, q2, q3 := quartiles(set)

	data = &quartilesItem{
		time:  set.Time,
		lo:    lo,
		q1:    int64(q1 + 0.5),
		q2:    int64(q2 + 0.5),
		q3:    int64(q3 + 0.5),
		hi:    hi,
		total: number,
	}

	return
}

// String returns string representation of the given quartilesItem.
func (self *quartilesItem) String() string {
	return fmt.Sprintf(
		"quartilesItem[time=%d, lo=%d, q1=%d, q2=%d, q3=%d, hi=%d, total=%d]",
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
func (*quartilesItem) rrdInfo() []string {
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
		"RRA:MAX:0.5:1:25920",       // 72 hours at 1 sample per 10 secs
		"RRA:MAX:0.5:60:4320",       // 1 month at 1 sample per 10 mins
		"RRA:MAX:0.5:2880:5475",     // 5 years at 1 sample per 8 hours
	}
}

// rrdTemplate returns template for RRDTool used to update data.
func (*quartilesItem) rrdTemplate() string {
	return "q1:q2:q3:lo:hi:total"
}

// rrdString returns a string matching template format with the data to
// update RRD files.
func (self *quartilesItem) rrdString() string {
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

// quartiles calculates quartiles for the given sample set.
func quartiles(set *types.SampleSet) (q1, q2, q3 float64) {
	number := int64(len(set.Values))
	q2index, q2 := median(set.Values)
	_, q1 = median(set.Values[:q2index+1])
	_, q3 = median(set.Values[number-q2index-1:])
	return
}

// median calculates value and index of the median for the given sample set.
func median(set []int) (index int64, median float64) {
	number := int64(len(set))
	var n float64 = float64(number-1) / 2.0
	k, d := math.Modf(n)
	index = int64(k)
	median = float64(set[index])
	if index+1 < number {
		median += d * float64(set[index+1]-set[index])
	}
	return
}

package writers

import (
	"fmt"
	"gorrdpd/types"
)

// Count writer is used to calculate positive and negative numbers.
type Count struct{}

// CountItem stores summary information about sample set.
type CountItem struct {
	// Timestamp of the sample set.
	time int64
	// Number of positive values.
	ok uint64
	// Number of negative values.
	fail uint64
}

// Name returns the name of the writer.
func (*Count) Name() string {
	return "count"
}

// Rollup performs summarization on the given sample set and writes
// results to RRD file.
func (self *Count) Rollup(set *types.SampleSet) {
	Rollup(self, set)
}

// BatchRollup performs summarization on the given list of sample sets and
// writes results to RRD files.
func (self *Count) BatchRollup(sets types.SampleSetsList) {
	BatchRollup(self, sets)
}

// rollupData performs summarization on the given sample set and returns
// CountItem with statistics.
func (self *Count) rollupData(set *types.SampleSet) (data dataItem) {
	var ok, fail uint64
	for _, elem := range set.Values {
		if elem > 0 {
			ok++
		} else if elem < 0 {
			fail++
		}
	}
	data = &CountItem{time: set.Time, ok: ok, fail: fail}
	return
}

// String returns string representation of the given CountItem.
func (self *CountItem) String() string {
	return fmt.Sprintf("CountItem[time=%d, ok=%d, fail=%d]", self.time, self.ok, self.fail)
}

// rrdInfo returns the list of parameters used to create RRD file.
func (*CountItem) rrdInfo() []string {
	return []string{
		"DS:ok:ABSOLUTE:600:0:U",
		"DS:fail:ABSOLUTE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",		// 72 hours at 1 sample per 10 secs
		"RRA:AVERAGE:0.5:60:4320",		// 1 month at 1 sample per 10 mins
		"RRA:AVERAGE:0.5:2880:5475",	// 5 years at 1 sample per 8 hours
	}
}

// rrdTemplate returns template for RRDTool used to update data.
func (*CountItem) rrdTemplate() string {
	return "ok:fail"
}

// rrdString returns a string matching template format with the data to
// update RRD files.
func (self *CountItem) rrdString() string {
	return fmt.Sprintf("%d:%d:%d", self.time, self.ok, self.fail)
}

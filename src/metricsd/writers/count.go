package writers

import (
	"fmt"
	"metricsd/types"
)

// Count writer is used to calculate positive and negative numbers.
type Count struct{
	*BaseWriter
}

// countItem stores summary information about sample set.
type countItem struct {
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

// rollupData performs summarization on the given sample set and returns
// countItem with statistics.
func (self *Count) rollupData(set *types.SampleSet) (data dataItem) {
	var ok, fail uint64
	for _, elem := range set.Values {
		if elem > 0 {
			ok++
		} else if elem < 0 {
			fail++
		}
	}
	data = &countItem{time: set.Time, ok: ok, fail: fail}
	return
}

// String returns string representation of the given countItem.
func (self *countItem) String() string {
	return fmt.Sprintf("countItem[time=%d, ok=%d, fail=%d]", self.time, self.ok, self.fail)
}

// rrdInfo returns the list of parameters used to create RRD file.
func (*countItem) rrdInfo() []string {
	return []string{
		"DS:ok:ABSOLUTE:600:0:U",
		"DS:fail:ABSOLUTE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",		// 72 hours at 1 sample per 10 secs
		"RRA:AVERAGE:0.5:60:4320",		// 1 month at 1 sample per 10 mins
		"RRA:AVERAGE:0.5:2880:5475",	// 5 years at 1 sample per 8 hours
	}
}

// rrdTemplate returns template for RRDTool used to update data.
func (*countItem) rrdTemplate() string {
	return "ok:fail"
}

// rrdString returns a string matching template format with the data to
// update RRD files.
func (self *countItem) rrdString() string {
	return fmt.Sprintf("%d:%d:%d", self.time, self.ok, self.fail)
}

package writers

import (
	"fmt"
	"math"
	"sort"
	"gorrdpd/types"
)

// Percentiles writer is used to calculate 90th and 95th percentiles, mean
// values, and standard deviations under percentiles.
type Percentiles struct{}

// PercentilesItem stores statistics information calculated by Percentiles
// writer.
type PercentilesItem struct {
	// Timestamp of the sample set.
	time int64
	// 90th percentile.
	pct90 int64
	// Mean value for metrics below the 90th percentile.
	pct90mean int64
	// Standard deviation for metrics below the 90th percentile.
	pct90dev int64
	// 95th percentile.
	pct95 int64
	// Mean value for metrics below the 95th percentile.
	pct95mean int64
	// Standard deviation for metrics below the 95th percentile.
	pct95dev int64
}

// Name returns the name of the writer.
func (self *Percentiles) Name() string {
	return "percentiles"
}

// Rollup performs summarization on the given sample set and writes
// results to RRD file.
func (self *Percentiles) Rollup(set *types.SampleSet) {
	Rollup(self, set)
}

// BatchRollup performs summarization on the given list of sample sets and
// writes results to RRD files.
func (self *Percentiles) BatchRollup(sets types.SampleSetsList) {
	BatchRollup(self, sets)
}

// rollupData performs summarization on the given sample set and returns
// PercentilesItem with statistics.
func (self *Percentiles) rollupData(set *types.SampleSet) (data dataItem) {
	if len(set.Values) == 0 {
		return
	}
	sort.Sort(set.Values)
	number := int64(len(set.Values))
	var n90 float64 = 90.0 / 100.0 * (float64(number) + 1)
	k90, d90 := math.Modf(n90)
	var n95 float64 = 95.0 / 100.0 * (float64(number) + 1)
	k95, d95 := math.Modf(n95)

	pct90index := int64(k90)
	pct95index := int64(k95)

	pct90 := float64(set.Values[pct90index - 1])
	pct95 := float64(set.Values[pct95index - 1])
	if pct90index > 1 && pct90index < number {
		pct90 += d90 * float64(set.Values[pct90index] - set.Values[pct90index - 1])
	}
	if pct95index > 1 && pct95index < number {
		pct95 += d95 * float64(set.Values[pct95index] - set.Values[pct95index - 1])
	}

	var pct90sum float64 = 0
	var pct95sum float64 = 0
	for idx, elem := range set.Values[0:pct95index] {
		if int64(idx) < pct90index {
			pct90sum += float64(elem)
		}
		pct95sum += float64(elem)
	}
	var pct90mean float64 = pct90sum / float64(pct90index)
	var pct95mean float64 = pct95sum / float64(pct95index)

	var pct90sqdiff float64 = 0
	var pct95sqdiff float64 = 0
	for idx, elem := range set.Values[0:pct95index] {
		if int64(idx) <= pct90index {
			pct90sqdiff += math.Pow(float64(pct90mean) - float64(elem), 2)
		}
		pct95sqdiff += math.Pow(float64(pct95mean) - float64(elem), 2)
	}

	data = &PercentilesItem{
		time:		set.Time,
		pct90:		int64(pct90 + 0.5),
		pct90mean:	int64(pct90mean + 0.5),
		pct90dev:	int64(math.Sqrt(pct90sqdiff / float64(pct90index)) + 0.5),
		pct95:		int64(pct95 + 0.5),
		pct95mean:	int64(pct95mean + 0.5),
		pct95dev:	int64(math.Sqrt(pct95sqdiff / float64(pct95index)) + 0.5),
	}
	return
}

// String returns string representation of the given PercentilesItem.
func (self *PercentilesItem) String() string {
	return fmt.Sprintf(
		"PercentilesItem[time=%d, pct90=%d, pct90mean=%d, pct90dev=%d, pct95=%d, pct95mean=%d, pct95dev=%d]",
		self.time,
		self.pct90,
		self.pct90mean,
		self.pct90dev,
		self.pct95,
		self.pct95mean,
		self.pct95dev,
	)
}

// rrdInfo returns the list of parameters used to create RRD file.
func (*PercentilesItem) rrdInfo() []string {
	return []string{
		"DS:pct90:GAUGE:600:0:U",
		"DS:pct90mean:GAUGE:600:0:U",
		"DS:pct90dev:GAUGE:600:0:U",
		"DS:pct95:GAUGE:600:0:U",
		"DS:pct95mean:GAUGE:600:0:U",
		"DS:pct95dev:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",		// 72 hours at 1 sample per 10 secs
		"RRA:AVERAGE:0.5:60:4320",		// 1 month at 1 sample per 10 mins
		"RRA:AVERAGE:0.5:2880:5475",	// 5 years at 1 sample per 8 hours
		"RRA:MAX:0.5:1:25920",			// 72 hours at 1 sample per 10 secs
		"RRA:MAX:0.5:60:4320",			// 1 month at 1 sample per 10 mins
		"RRA:MAX:0.5:2880:5475",		// 5 years at 1 sample per 8 hours
	}
}

// rrdTemplate returns template for RRDTool used to update data.
func (*PercentilesItem) rrdTemplate() string {
	return "pct90:pct90mean:pct90dev:pct95:pct95mean:pct95dev"
}

// rrdString returns a string matching template format with the data to
// update RRD files.
func (self *PercentilesItem) rrdString() string {
	return fmt.Sprintf(
		"%d:%d:%d:%d:%d:%d:%d",
		self.time,
		self.pct90,
		self.pct90mean,
		self.pct90dev,
		self.pct95,
		self.pct95mean,
		self.pct90dev,
	)
}

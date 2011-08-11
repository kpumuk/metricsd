package writers

import (
	. "launchpad.net/gocheck"
	"testing"
	"metricsd/types"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

func createSampleSet(time int64, values ...int) (ss *types.SampleSet) {
	ss = types.NewSampleSet(time, "src", "metric")
	fillSampleSet(ss, values...)
	return
}

func fillSampleSet(ss *types.SampleSet, values ...int) {
	for _, value := range values {
		ss.Add(value)
	}
}

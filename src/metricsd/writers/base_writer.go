package writers

import (
	"metricsd/types"
)

type Writer interface {
	Name() string
	Rollup(set *types.SampleSet)
	BatchRollup(sets []*types.SampleSet)
	// Private methods
	rollupData(set *types.SampleSet) dataItem
}

type BaseWriter struct {}

// Rollup performs summarization on the given sample set and writes
// results to RRD file.
func (writer *BaseWriter) Rollup(set *types.SampleSet) {
	Rollup(writer, set)
}

// BatchRollup performs summarization on the given list of sample sets and
// writes results to RRD files.
func (writer *BaseWriter) BatchRollup(sets []*types.SampleSet) {
	BatchRollup(writer, sets)
}

func (writer *BaseWriter) Name() string {
	panic("You should implement Name() func")
}

func (writer *BaseWriter) rollupData(set *types.SampleSet) dataItem {
	panic("You should implement rollupData() func")
}

type dataItem interface {
	rrdInfo() []string
	rrdTemplate() string
	rrdString() string
	String() string
}

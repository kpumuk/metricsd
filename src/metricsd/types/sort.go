package types

import (
	"sort"
)

// SampleSetSlice attaches the methods of sort.Interface to []*SampleSet, sorting in increasing order.
type SampleSetSlice []*SampleSet

func (p SampleSetSlice) Len() int           { return len(p) }
func (p SampleSetSlice) Less(i, j int) bool { return p[i].Less(p[j]) }
func (p SampleSetSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// SliceSlice attaches the methods of sort.Interface to []*Slice, sorting in increasing order.
type SliceSlice []*Slice

func (p SliceSlice) Len() int           { return len(p) }
func (p SliceSlice) Less(i, j int) bool { return p[i].Less(p[j]) }
func (p SliceSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// SortSampleSets sorts a slice of *SampleSet in increasing order.
func SortSampleSets(a []*SampleSet) { sort.Sort(SampleSetSlice(a)) }
// SliceSlices sorts a slice of *Slice in increasing order.
func SortSlices(a []*Slice) { sort.Sort(SliceSlice(a)) }

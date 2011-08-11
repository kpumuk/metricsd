package types

import (
	"fmt"
)

type Slice struct {
	Time int64
	Sets map[string]*SampleSet
}

func NewSlice(time int64) *Slice {
	return &Slice{
		Time: time,
		Sets: make(map[string]*SampleSet),
	}
}

func (slice *Slice) Less(sliceToCompare *Slice) bool {
	return slice.Time < sliceToCompare.Time
}

func (slice *Slice) Add(event *Event) {
	slice.getSampleSet(slice.getAllSampleSetKey(event.Name), "all", event.Name).Add(event.Value)
	if event.Source != "all" {
		slice.getSampleSet(slice.getMachineSampleSetKey(event.Source, event.Name), event.Source, event.Name).Add(event.Value)
	}
}

func (slice *Slice) String() string {
	return fmt.Sprintf(
		"Slice[time=%d, size=%d]",
		slice.Time,
		len(slice.Sets),
	)
}

func (slice *Slice) getSampleSet(key, source, name string) *SampleSet {
	if _, found := slice.Sets[key]; !found {
		slice.Sets[key] = NewSampleSet(slice.Time, source, name)
	}
	return slice.Sets[key]
}

func (slice *Slice) getAllSampleSetKey(name string) string {
	return "all-" + name
}

func (slice *Slice) getMachineSampleSetKey(source, name string) string {
	return source + "-" + name
}

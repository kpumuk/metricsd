package types

import (
	"fmt"
)

type SampleSet struct {
	Time   int64
	Source string
	Name   string
	Values []int
}

func NewSampleSet(time int64, source, name string) *SampleSet {
	return &SampleSet{
		Time:   time,
		Source: source,
		Name:   name,
		Values: make([]int, 0, 8),
	}
}

func (set *SampleSet) Add(event *Event) {
	set.Values = append(set.Values, event.Value)
}

func (set *SampleSet) Less(setToCompare *SampleSet) bool {
	return set.Source < setToCompare.Source ||
		(set.Source == setToCompare.Source && set.Name < setToCompare.Name) ||
		(set.Source == setToCompare.Source && set.Name == setToCompare.Name && set.Time < setToCompare.Time)
}

func (set *SampleSet) String() string {
	return fmt.Sprintf(
		"SampleSet[source=%s, name=%s, time=%d, size=%d]",
		set.Source,
		set.Name,
		set.Time,
		len(set.Values),
	)
}

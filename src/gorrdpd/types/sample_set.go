package types

import (
    "fmt"
    "container/vector"
)

type SampleSet struct {
    Time   int64
    Source string
    Name   string
    Values *vector.IntVector
}

func NewSampleSet(time int64, source, name string) *SampleSet {
    return &SampleSet{
        Time:   time,
        Source: source,
        Name:   name,
        Values: new(vector.IntVector),
    }
}

func (set *SampleSet) Add(message *Message) {
    set.Values.Push(message.Value)
}

func (set *SampleSet) Less(setToCompare interface{}) bool {
    s := setToCompare.(*SampleSet)
    return set.Source < s.Source ||
        (set.Source == s.Source && set.Name < s.Name) ||
        (set.Source == s.Source && set.Name == s.Name && set.Time < s.Time)
}

func (set *SampleSet) String() string {
    return fmt.Sprintf(
        "SampleSet[source=%s, name=%s, time=%d, size=%d]",
        set.Source,
        set.Name,
        set.Time,
        set.Values.Len(),
    )
}

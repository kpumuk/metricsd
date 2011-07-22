package types

import (
    "fmt"
)

type Slice struct {
    Time int64
    Sets map[string]*SampleSet
}

func NewSlice(time int64) *Slice {
    return &Slice{Time: time, Sets: make(map[string]*SampleSet)}
}

func (slice *Slice) Less(sliceToCompare *Slice) bool {
    return slice.Time < sliceToCompare.Time
}

func (slice *Slice) Add(message *Message) {
    slice.getSampleSet(slice.getAllSampleSetKey(message), "all", message.Name).Add(message)
    if message.Source != "all" {
        slice.getSampleSet(slice.getMachineSampleSetKey(message), message.Source, message.Name).Add(message)
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

func (slice *Slice) getAllSampleSetKey(message *Message) string {
    return fmt.Sprintf("all-%s", message.Name)
}

func (slice *Slice) getMachineSampleSetKey(message *Message) string {
    return fmt.Sprintf("%s-%s", message.Source, message.Name)
}

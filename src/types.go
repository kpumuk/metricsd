package types

import (
    "container/vector"
    "fmt"
    "sort"
    "time"
)

/******************************************************************************/

type Message struct {
    Source string
    Name   string
    Value  int
}

func NewMessage(source string, name string, value int) *Message {
    return &Message { Source: source, Name: name, Value: value };
}

func (message *Message) String() string {
    return fmt.Sprintf("Message[source=%s, name=%s, value=%d]", message.Source, message.Name, message.Value)
}

/******************************************************************************/

type Slices struct {
    Interval int64
    Slices map[int64] *Slice
}

func NewSlices(sliceInterval int) *Slices {
    return &Slices { Slices: make(map[int64] *Slice), Interval: int64(sliceInterval) }
}

func (slices *Slices) Add(message *Message) {
    slices.getCurrentSlice().Add(message)
}

func (slices *Slices) ExtractClosedSlices(force bool) *vector.Vector {
    current := slices.getCurrentSliceNumber()
    closedSlices := new(vector.Vector)
    for number, slice := range slices.Slices {
        if number < current || force {
            closedSlices.Push(slice)
            slices.Slices[number] = nil, false
        }
    }
    sort.Sort(closedSlices)
    return closedSlices
}

func (slices *Slices) String() string {
    return fmt.Sprintf("Slices[interval=%d, size=%d]", slices.Interval, len(slices.Slices))
}

func (slices *Slices) getCurrentSlice() *Slice {
    number := slices.getCurrentSliceNumber()
    if _, found := slices.Slices[number]; !found {
        slices.Slices[number] = NewSlice(number * slices.Interval)
    }
    return slices.Slices[number]
}

func (slices *Slices) getCurrentSliceNumber() int64 {
    return time.Seconds() / slices.Interval
}

/******************************************************************************/

type Slice struct {
    Time int64
    Sets map[string] *SampleSet
}

func NewSlice(time int64) *Slice {
    return &Slice { Time: time, Sets: make(map[string] *SampleSet) }
}

func (slice *Slice) Less(sliceToCompare interface{}) bool {
    return slice.Time < sliceToCompare.(*Slice).Time
}

func (slice *Slice) Add(message *Message) {
    slice.getSampleSet(slice.getAllSampleSetKey(message)).Add(message)
    slice.getSampleSet(slice.getMachineSampleSetKey(message)).Add(message)
}

func (slice *Slice) String() string {
    return fmt.Sprintf("Slice[time=%d, size=%d]", slice.Time, len(slice.Sets))
}

func (slice *Slice) getSampleSet(key string) *SampleSet {
    if _, found := slice.Sets[key]; !found {
        slice.Sets[key] = NewSampleSet(slice.Time, key)
    }
    return slice.Sets[key]
}

func (slice *Slice) getAllSampleSetKey(message *Message) string {
    return fmt.Sprintf("all-%s", message.Name)
}

func (slice *Slice) getMachineSampleSetKey(message *Message) string {
    return fmt.Sprintf("%s-%s", message.Source, message.Name)
}

/******************************************************************************/

type SampleSet struct {
    Time   int64
    Key    string
    Values *vector.IntVector
}

func NewSampleSet(time int64, key string) *SampleSet {
    return &SampleSet { Time: time, Key: key, Values: new(vector.IntVector) }
}

func (set *SampleSet) Add(message *Message) {
    set.Values.Push(message.Value)
}

func (set *SampleSet) String() string {
    return fmt.Sprintf("SampleSet[time=%d, key=%s, size=%d]", set.Time, set.Key, set.Values.Len())
}

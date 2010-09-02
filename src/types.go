package types

import (
    "container/vector"
    "fmt"
    "sort"
    "time"
)

/***** Message ****************************************************************/

// A Message contains information about message
type Message struct {
    Source string   // message source (IP address, DNS name, or custom string)
    Name   string   // metric's name
    Value  int      // metric's value
}

// NewMessage creates a new message instance.
func NewMessage(source string, name string, value int) *Message {
    return &Message { Source: source, Name: name, Value: value };
}

// String converts an instance of message struct to string.
func (message *Message) String() string {
    return fmt.Sprintf(
        "Message[source=%s, name=%s, value=%d]",
        message.Source,
        message.Name,
        message.Value,
    )
}

/******************************************************************************/

type Slices struct {
    Interval int64
    Slices map[int64] *Slice
}

func NewSlices(sliceInterval int) *Slices {
    return &Slices {
        Slices:   make(map[int64] *Slice),
        Interval: int64(sliceInterval),
    }
}

func (slices *Slices) Add(message *Message) {
    slices.getCurrentSlice().Add(message)
}

func (slices *Slices) ExtractClosedSlices(force bool) (closedSlices *vector.Vector) {
    var current int64
    if force { current = -1 } else { current = slices.getCurrentSliceNumber() }

    // Calculate total number of closed slices (to avoid vector reallocs)
    totalClosedSlices := 0
    slices.eachClosedSlice(current, func(number int64, slice *Slice) {
        totalClosedSlices += 1
    })

    // Create an array to store slices
    closedSlices = new(vector.Vector).Resize(0, totalClosedSlices)
    slices.eachClosedSlice(current, func(number int64, slice *Slice) {
        closedSlices.Push(slice)
        slices.Slices[number] = nil, false
    })
    sort.Sort(closedSlices)
    return
}

// ExtractClosedSampleSets finds closed slices, and stores all sample sets from them
// in an array. Processed slices will be removed from the list of active slices.
func (slices *Slices) ExtractClosedSampleSets(force bool) (closedSampleSets *vector.Vector) {
    var current int64
    if force { current = -1 } else { current = slices.getCurrentSliceNumber() }

    // Calculate total number of closed sample sets (to avoid vector reallocs)
    totalSampleSets := 0
    slices.eachClosedSlice(current, func(number int64, slice *Slice) {
        totalSampleSets += len(slice.Sets)
    })

    // Create an array to store sample sets
    closedSampleSets = new(vector.Vector).Resize(0, totalSampleSets)
    slices.eachClosedSlice(current, func(number int64, slice *Slice) {
        for _, set := range slice.Sets { closedSampleSets.Push(set) }
        slices.Slices[number] = nil, false
    })
    sort.Sort(closedSampleSets)
    return
}

func (slices *Slices) String() string {
    return fmt.Sprintf(
        "Slices[interval=%d, size=%d]",
        slices.Interval,
        len(slices.Slices),
    )
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

func (slices *Slices) eachClosedSlice(current int64, f func(number int64, slice *Slice)) {
    for number, slice := range slices.Slices {
        if number < current || current < 0 {
            f(number, slice)
        }
    }
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
        slice.Sets[key] = NewSampleSet(slice.Time, key, source, name)
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
    Source string
    Name   string
    Values *vector.IntVector
}

func NewSampleSet(time int64, key, source, name string) *SampleSet {
    return &SampleSet {
        Time:   time,
        Key:    key,
        Source: source,
        Name:   name,
        Values: new(vector.IntVector),
    }
}

func (set *SampleSet) Add(message *Message) {
    set.Values.Push(message.Value)
}

func (set *SampleSet) Less(setToCompare interface{}) bool {
    return set.Key < setToCompare.(*SampleSet).Key ||
        (set.Key == setToCompare.(*SampleSet).Key && set.Time < setToCompare.(*SampleSet).Time)
}

func (set *SampleSet) String() string {
    return fmt.Sprintf(
        "SampleSet[time=%d, key=%s, size=%d]",
        set.Time,
        set.Key,
        set.Values.Len(),
    )
}

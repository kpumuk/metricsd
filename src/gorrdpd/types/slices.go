package types

import (
    "container/vector"
    "fmt"
    "sort"
	"time"
)

type Slices struct {
    Interval int64
    Slices   map[int64]*Slice
}

func NewSlices(sliceInterval int) *Slices {
    return &Slices{
        Slices:   make(map[int64]*Slice),
        Interval: int64(sliceInterval),
    }
}

func (slices *Slices) Add(message *Message) {
    slices.getCurrentSlice().Add(message)
}

func (slices *Slices) ExtractClosedSlices(force bool) (closedSlices *vector.Vector) {
    var current int64
    if force {
        current = -1
    } else {
        current = slices.getCurrentSliceNumber()
    }

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
    if force {
        current = -1
    } else {
        current = slices.getCurrentSliceNumber()
    }

    // Calculate total number of closed sample sets (to avoid vector reallocs)
    totalSampleSets := 0
    slices.eachClosedSlice(current, func(number int64, slice *Slice) {
        totalSampleSets += len(slice.Sets)
    })

    // Create an array to store sample sets
    closedSampleSets = new(vector.Vector).Resize(0, totalSampleSets)
    slices.eachClosedSlice(current, func(number int64, slice *Slice) {
        for _, set := range slice.Sets {
            closedSampleSets.Push(set)
        }
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

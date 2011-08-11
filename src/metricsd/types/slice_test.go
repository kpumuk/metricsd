package types

import (
	. "launchpad.net/gocheck"
	"testing"
)

type SliceS struct {
	slice *Slice
}

var _ = Suite(&SliceS{})


func (s *SliceS) SetUpTest(c *C) {
	s.slice = NewSlice(10)
}

func (s *SliceS) TestGetAllSampleSetKey(c *C) {
	key := s.slice.getSampleSetKey("all", "metric")
	c.Check(key, Equals, "all-metric")
}

func (s *SliceS) TestGetMachineSampleSetKey(c *C) {
	key := s.slice.getSampleSetKey("src", "metric")
	c.Check(key, Equals, "src-metric")
}

func BenchmarkSliceAdd(b *testing.B) {
	b.StopTimer()
	ss := NewSlice(10)
	evt := &Event{Source: "src", Name: "metric", Value: 10}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ss.Add(evt)
	}

	b.StopTimer()
}

func BenchmarkSliceGetSampleSetKey(b *testing.B) {
	b.StopTimer()
	ss := NewSlice(10)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ss.getSampleSetKey("src", "metric")
	}

	b.StopTimer()
}

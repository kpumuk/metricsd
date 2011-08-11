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
	key := s.slice.getAllSampleSetKey("metric")
	c.Check(key, Equals, "all-metric")
}

func (s *SliceS) TestGetMachineSampleSetKey(c *C) {
	key := s.slice.getMachineSampleSetKey("src", "metric")
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

func BenchmarkSliceGetAllSampleSetKey(b *testing.B) {
	b.StopTimer()
	ss := NewSlice(10)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ss.getAllSampleSetKey("metric")
	}

	b.StopTimer()
}

func BenchmarkSliceGetMachineSampleSetKey(b *testing.B) {
	b.StopTimer()
	ss := NewSlice(10)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ss.getMachineSampleSetKey("src", "metric")
	}

	b.StopTimer()
}

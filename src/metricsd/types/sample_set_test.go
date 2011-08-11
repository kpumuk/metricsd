package types

import (
	. "launchpad.net/gocheck"
	"testing"
)

type SampleSetS struct{}

var _ = Suite(&SampleSetS{})

func BenchmarkSampleSetAdd(b *testing.B) {
	b.StopTimer()
	ss := NewSampleSet(10, "src", "metric")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ss.Add(i)
	}

	b.StopTimer()
}

package writers

import (
	. "launchpad.net/gocheck"
	"testing"
	"gorrdpd/types"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type S struct {
	count *Count
}
var _ = Suite(&S{})

func (s *S) SetUpTest(c *C) {
    s.count = &Count{}
}

func (s *S) TestRollupDataWithEmptySampleSet(c *C) {
	ss := types.NewSampleSet(10, "src", "metric")
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &CountItem{time: 10, ok: 0, fail: 0})
}

func (s *S) TestRollupDataWithSampleSetContainsZero(c *C) {
	ss := types.NewSampleSet(11, "src", "metric")
	ss.Add(&types.Event{Value:0})
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &CountItem{time: 11, ok: 0, fail: 0})
}

func (s *S) TestRollupDataWithSimpleSampleSet(c *C) {
	ss := types.NewSampleSet(12, "src", "metric")
	ss.Add(&types.Event{Value:1})
	ss.Add(&types.Event{Value:-1})
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &CountItem{time: 12, ok: 1, fail: 1})
}

func (s *S) TestRollupDataWithComplexSampleSet(c *C) {
	ss := types.NewSampleSet(13, "src", "metric")
	ss.Add(&types.Event{Value:5})
	ss.Add(&types.Event{Value:1})
	ss.Add(&types.Event{Value:1000})
	ss.Add(&types.Event{Value:-5})
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &CountItem{time: 13, ok: 3, fail: 1})
}

package writers

import (
	. "launchpad.net/gocheck"
	"gorrdpd/types"
)

type CountS struct {
	count *Count
}
var _ = Suite(&CountS{})

func (s *CountS) SetUpTest(c *C) {
    s.count = &Count{}
}

func (s *CountS) TestRollupDataWithEmptySampleSet(c *C) {
	ss := types.NewSampleSet(10, "src", "metric")
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 10, ok: 0, fail: 0})
}

func (s *CountS) TestRollupDataWithSampleSetContainsZero(c *C) {
	ss := types.NewSampleSet(11, "src", "metric")
	ss.Add(&types.Event{Value:0})
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 11, ok: 0, fail: 0})
}

func (s *CountS) TestRollupDataWithSimpleSampleSet(c *C) {
	ss := types.NewSampleSet(12, "src", "metric")
	ss.Add(&types.Event{Value:1})
	ss.Add(&types.Event{Value:-1})
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 12, ok: 1, fail: 1})
}

func (s *CountS) TestRollupDataWithComplexSampleSet(c *C) {
	ss := types.NewSampleSet(13, "src", "metric")
	ss.Add(&types.Event{Value:5})
	ss.Add(&types.Event{Value:1})
	ss.Add(&types.Event{Value:1000})
	ss.Add(&types.Event{Value:-5})
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 13, ok: 3, fail: 1})
}

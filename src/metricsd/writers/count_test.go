package writers

import (
	. "launchpad.net/gocheck"
)

type CountS struct {
	count *Count
}

var _ = Suite(&CountS{})

func (s *CountS) SetUpTest(c *C) {
	s.count = &Count{}
}

func (s *CountS) TestRollupDataWithEmptySampleSet(c *C) {
	ss := createSampleSet(1000)
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 1000, ok: 0, fail: 0})
}

func (s *CountS) TestRollupDataWithSampleSetContainsZero(c *C) {
	ss := createSampleSet(2000, 0)
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 2000, ok: 0, fail: 0})
}

func (s *CountS) TestRollupDataWithSimpleSampleSet(c *C) {
	ss := createSampleSet(3000, 1, -1)
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 3000, ok: 1, fail: 1})
}

func (s *CountS) TestRollupDataWithComplexSampleSet(c *C) {
	ss := createSampleSet(4000, 5, 1, 1000, -5)
	data := s.count.rollupData(ss)
	c.Check(data, Equals, &countItem{time: 4000, ok: 3, fail: 1})
}

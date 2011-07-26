package writers

import (
	. "launchpad.net/gocheck"
)

type QuartilesS struct {
	quartiles *Quartiles
}
var _ = Suite(&QuartilesS{})

func (s *QuartilesS) SetUpTest(c *C) {
    s.quartiles = &Quartiles{}
}

func (s *QuartilesS) TestRollupDataWithEmptySampleSet(c *C) {
	ss := createSampleSet(1000)
	data := s.quartiles.rollupData(ss)
	c.Check(data, IsNil)
}

func (s *QuartilesS) TestRollupDataWithSampleSetWith1Item(c *C) {
	ss := createSampleSet(2000, 10)
	data := s.quartiles.rollupData(ss)
	c.Check(data, Equals, &quartilesItem{time: 2000, lo: 10, q1: 10, q2: 10, q3: 10, hi: 10, total: 1})
}

func (s *QuartilesS) TestRollupDataWithSampleSetWith2Items(c *C) {
	ss := createSampleSet(3000, 10, 20)
	data := s.quartiles.rollupData(ss)
	c.Check(data, Equals, &quartilesItem{time: 3000, lo: 10, q1: 10, q2: 15, q3: 20, hi: 20, total: 2})
}

func (s *QuartilesS) TestRollupDataWithSampleSetWith3Items(c *C) {
	ss := createSampleSet(4000, 10, 20, 30)
	data := s.quartiles.rollupData(ss)
	c.Check(data, Equals, &quartilesItem{time: 4000, lo: 10, q1: 15, q2: 20, q3: 25, hi: 30, total: 3})
}

func (s *QuartilesS) TestRollupDataWithSampleSetWith5Items(c *C) {
	ss := createSampleSet(5000, 36, 7, 15, 40, 41, 39)
	data := s.quartiles.rollupData(ss)
	c.Check(data, Equals, &quartilesItem{time: 5000, lo: 7, q1: 15, q2: 38, q3: 40, hi: 41, total: 6})
}

func (s *QuartilesS) TestRollupDataWithLargeSampleSet(c *C) {
	ss := createSampleSet(6000, 6, 47, 49, 15, 42, 41, 7, 39, 43, 40, 36)
	data := s.quartiles.rollupData(ss)
	c.Check(data, Equals, &quartilesItem{time: 6000, lo: 6, q1: 26, q2: 40, q3: 43, hi: 49, total: 11})
}

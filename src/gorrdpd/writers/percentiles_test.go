package writers

import (
	. "launchpad.net/gocheck"
	"gorrdpd/types"
)

type PercentilesS struct {
	percentiles *Percentiles
}
var _ = Suite(&PercentilesS{})

func (s *PercentilesS) SetUpTest(c *C) {
    s.percentiles = &Percentiles{}
}

func (s *PercentilesS) TestRollupDataWithEmptySampleSet(c *C) {
	ss := createSampleSet(1000)
	data := s.percentiles.rollupData(ss)
	c.Check(data, IsNil)
}

func (s *PercentilesS) TestRollupDataWithSampleSetWith1Item(c *C) {
	ss := createSampleSet(2000, 10)
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 2000, pct90: 10, pct90mean: 10, pct90dev: 0, pct95: 10, pct95mean: 10, pct95dev: 0})
}

func (s *PercentilesS) TestRollupDataWithSampleSetWith2Items(c *C) {
	ss := createSampleSet(3000, 10, 20)
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 3000, pct90: 20, pct90mean: 15, pct90dev: 5, pct95: 20, pct95mean: 15, pct95dev: 5})
}

func (s *PercentilesS) TestRollupDataWithSampleSetWith3Items(c *C) {
	ss := createSampleSet(4000, 10, 20, 30)
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 4000, pct90: 30, pct90mean: 20, pct90dev: 8, pct95: 30, pct95mean: 20, pct95dev: 8})
}

func (s *PercentilesS) TestRollupDataWithSimpleSampleSet(c *C) {
	ss := createSampleSet(5000, 15, 20, 35, 40, 50)
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 5000, pct90: 50, pct90mean: 32, pct90dev: 13, pct95: 50, pct95mean: 32, pct95dev: 13})
}

func (s *PercentilesS) TestRollupDataWithComplexSampleSet(c *C) {
	ss := createSampleSet(6000)
	for i := 1; i < 100; i++ {
		ss.Add(&types.Event{Value:i * 10})
	}
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 6000, pct90: 900, pct90mean: 455, pct90dev: 264, pct95: 950, pct95mean: 480, pct95dev: 274})
}
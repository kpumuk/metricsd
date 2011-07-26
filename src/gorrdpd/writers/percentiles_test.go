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
	ss := types.NewSampleSet(10, "src", "metric")
	data := s.percentiles.rollupData(ss)
	c.Check(data, IsNil)
}

func (s *PercentilesS) TestRollupDataWithSampleSetWith1Item(c *C) {
	ss := types.NewSampleSet(11, "src", "metric")
	ss.Add(&types.Event{Value:10})
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 11, pct90: 10, pct90mean: 10, pct90dev: 0, pct95: 10, pct95mean: 10, pct95dev: 0})
}

func (s *PercentilesS) TestRollupDataWithSampleSetWith2Items(c *C) {
	ss := types.NewSampleSet(12, "src", "metric")
	ss.Add(&types.Event{Value:10})
	ss.Add(&types.Event{Value:20})
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 12, pct90: 20, pct90mean: 15, pct90dev: 5, pct95: 20, pct95mean: 15, pct95dev: 5})
}

func (s *PercentilesS) TestRollupDataWithSampleSetWith3Items(c *C) {
	ss := types.NewSampleSet(13, "src", "metric")
	ss.Add(&types.Event{Value:10})
	ss.Add(&types.Event{Value:20})
	ss.Add(&types.Event{Value:30})
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 13, pct90: 30, pct90mean: 20, pct90dev: 8, pct95: 30, pct95mean: 20, pct95dev: 8})
}

func (s *PercentilesS) TestRollupDataWithSimpleSampleSet(c *C) {
	ss := types.NewSampleSet(14, "src", "metric")
	ss.Add(&types.Event{Value:15})
	ss.Add(&types.Event{Value:20})
	ss.Add(&types.Event{Value:35})
	ss.Add(&types.Event{Value:40})
	ss.Add(&types.Event{Value:50})
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 14, pct90: 50, pct90mean: 32, pct90dev: 13, pct95: 50, pct95mean: 32, pct95dev: 13})
}

func (s *PercentilesS) TestRollupDataWithComplexSampleSet(c *C) {
	ss := types.NewSampleSet(14, "src", "metric")
	for i := 1; i < 100; i++ {
		ss.Add(&types.Event{Value:i * 10})
	}
	data := s.percentiles.rollupData(ss)
	c.Check(data, Equals, &percentilesItem{time: 14, pct90: 900, pct90mean: 455, pct90dev: 264, pct95: 950, pct95mean: 480, pct95dev: 274})
}
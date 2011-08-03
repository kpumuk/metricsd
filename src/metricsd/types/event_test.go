package types

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }


type S struct{}

var _ = Suite(&S{})

func (s *S) TestNewEvent(c *C) {
	event := NewEvent("src", "msg", 10)
	c.Check(event.Source, Equals, "src")
	c.Check(event.Name, Equals, "msg")
	c.Check(event.Value, Equals, 10)
}

func (s *S) TestEventString(c *C) {
	event := NewEvent("src", "msg", 10)
	c.Check(event.String(), Equals, "Event[source=src, name=msg, value=10]")
}

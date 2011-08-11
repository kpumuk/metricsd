package types

import (
	. "launchpad.net/gocheck"
)

type EventS struct{}

var _ = Suite(&EventS{})

func (s *EventS) TestNewEvent(c *C) {
	event := NewEvent("src", "msg", 10)
	c.Check(event.Source, Equals, "src")
	c.Check(event.Name, Equals, "msg")
	c.Check(event.Value, Equals, 10)
}

func (s *EventS) TestEventString(c *C) {
	event := NewEvent("src", "msg", 10)
	c.Check(event.String(), Equals, "Event[source=src, name=msg, value=10]")
}

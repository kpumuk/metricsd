package types

import (
	"fmt"
)

type MetricValue int

// A Event contains information about the event.
type Event struct {
	Source string // event source (IP address, DNS name, or custom string)
	Name   string // metric's name
	Value  int	// metric's value
}

// NewEvent returns a new Event with the given source, name, and value.
func NewEvent(source string, name string, value int) *Event {
	return &Event{Source: source, Name: name, Value: value}
}

// String converts an instance of event struct to string.
func (event *Event) String() string {
	if event == nil {
		return "Event[nil]"
	}
	return fmt.Sprintf(
		"Event[source=%s, name=%s, value=%d]",
		event.Source,
		event.Name,
		event.Value,
	)
}

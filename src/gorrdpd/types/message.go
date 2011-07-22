package types

import (
    "fmt"
)

// A Message contains information about message
type Message struct {
    Source string // message source (IP address, DNS name, or custom string)
    Name   string // metric's name
    Value  int    // metric's value
}

// NewMessage creates a new message instance.
func NewMessage(source string, name string, value int) *Message {
    return &Message{Source: source, Name: name, Value: value}
}

// String converts an instance of message struct to string.
func (message *Message) String() string {
    if message == nil {
        return "Message[nil]"
    }
    return fmt.Sprintf(
        "Message[source=%s, name=%s, value=%d]",
        message.Source,
        message.Name,
        message.Value,
    )
}

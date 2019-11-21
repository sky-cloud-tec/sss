package common

import "time"

// Message is a log message, with a reception timestamp and sequence number.
type Message struct {
	Text          string                 // Delimited log line
	Parsed        map[string]interface{} // If non-nil, contains parsed fields
	ReceptionTime time.Time              // Time log line was received
	SourceIP      string                 // Sender's IP address
}

// NewMessage returns a new Message.
func NewMessage() *Message {
	return &Message{}
}

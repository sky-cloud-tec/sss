// Simple syslog server.
// Copyright (C) 2019  sky-cloud.net
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

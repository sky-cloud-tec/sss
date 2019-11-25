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

package syslog

import (
	"regexp"
	"strings"
)

const (
	// RFC5424DelimiterPrefix indicates the start of a syslog line
	RFC5424DelimiterPrefix = `<[0-9]{1,3}>[0-9]\s`
)

var syslogRegex *regexp.Regexp
var startRegex *regexp.Regexp
var runRegex *regexp.Regexp

func init() {
	syslogRegex = regexp.MustCompile(RFC5424DelimiterPrefix)
	startRegex = regexp.MustCompile(RFC5424DelimiterPrefix + `$`)
	runRegex = regexp.MustCompile(`\n` + RFC5424DelimiterPrefix)
}

// A RFC5424Delimiter detects when Syslog lines start.
type RFC5424Delimiter struct {
	buffer []byte
	regex  *regexp.Regexp
}

// NewRFC5424Delimiter returns an initialized RFC5424Delimiter.
func NewRFC5424Delimiter(maxSize int) *RFC5424Delimiter {
	s := &RFC5424Delimiter{}
	s.buffer = make([]byte, 0, maxSize)
	s.regex = startRegex
	return s
}

// Push a byte into the RFC5424Delimiter. If the byte results in a
// a new Syslog message, it'll be flagged via the bool.
func (s *RFC5424Delimiter) Push(b byte) (string, bool) {
	s.buffer = append(s.buffer, b)
	delimiter := s.regex.FindIndex(s.buffer)
	if delimiter == nil {
		return "", false
	}

	if s.regex == startRegex {
		// First match -- switch to the regex for embedded lines, and
		// drop any leading characters.
		s.buffer = s.buffer[delimiter[0]:]
		s.regex = runRegex
		return "", false
	}

	dispatch := strings.TrimRight(string(s.buffer[:delimiter[0]]), "\r")
	s.buffer = s.buffer[delimiter[0]+1:]
	return dispatch, true
}

// Vestige returns the bytes which have been pushed to RFC5424Delimiter, since
// the last Syslog message was returned, but only if the buffer appears
// to be a valid syslog message.
func (s *RFC5424Delimiter) Vestige() (string, bool) {
	delimiter := syslogRegex.FindIndex(s.buffer)
	if delimiter == nil {
		s.buffer = nil
		return "", false
	}
	dispatch := strings.TrimRight(string(s.buffer), "\r\n")
	s.buffer = nil
	return dispatch, true
}

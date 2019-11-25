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
	// RFC3164DelimiterPrefix indicates the start of a syslog line
	RFC3164DelimiterPrefix = `<[0-9]{1,3}>`
)

var rfc3164Regex *regexp.Regexp
var rfc3164startRegex *regexp.Regexp
var rfc3164runRegex *regexp.Regexp

func init() {
	rfc3164Regex = regexp.MustCompile(RFC3164DelimiterPrefix)
	rfc3164startRegex = regexp.MustCompile(RFC3164DelimiterPrefix + `$`)
	rfc3164runRegex = regexp.MustCompile(`\n` + RFC3164DelimiterPrefix)
}

// A RFC3164Delimiter detects when Syslog lines start.
type RFC3164Delimiter struct {
	buffer []byte
	regex  *regexp.Regexp
}

// NewRFC3164Delimiter returns an initialized RFC3164Delimiter.
func NewRFC3164Delimiter(maxSize int) *RFC3164Delimiter {
	s := &RFC3164Delimiter{}
	s.buffer = make([]byte, 0, maxSize)
	s.regex = rfc3164startRegex
	return s
}

// Push a byte into the RFC3164Delimiter. If the byte results in a
// a new RFC3164 message, it'll be flagged via the bool.
func (s *RFC3164Delimiter) Push(b byte) (string, bool) {
	s.buffer = append(s.buffer, b)
	delimiter := s.regex.FindIndex(s.buffer)
	if delimiter == nil {
		return "", false
	}

	if s.regex == rfc3164startRegex {
		// First match -- switch to the regex for embedded lines, and
		// drop any leading characters.
		s.buffer = s.buffer[delimiter[0]:]
		s.regex = rfc3164runRegex
		return "", false
	}

	dispatch := strings.TrimRight(string(s.buffer[:delimiter[0]]), "\r")
	s.buffer = s.buffer[delimiter[0]+1:]
	return dispatch, true
}

// Vestige returns the bytes which have been pushed to RFC3164Delimiter, since
// the last RFC3164 message was returned, but only if the buffer appears
// to be a valid syslog message.
func (s *RFC3164Delimiter) Vestige() (string, bool) {
	delimiter := rfc3164Regex.FindIndex(s.buffer)
	if delimiter == nil {
		s.buffer = nil
		return "", false
	}
	dispatch := strings.TrimRight(string(s.buffer), "\r\n")
	s.buffer = nil
	return dispatch, true
}

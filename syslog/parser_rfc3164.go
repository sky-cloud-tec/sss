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
	"strconv"
)

// RFC3164 represents a parser for RFC3164-compliant log messages
// BUT made some modifications to make it compatible with juniper syslog messages
type RFC3164 struct {
	matcher *regexp.Regexp
}

func (p *Parser) newRFC3164Parser() {
	p.rfc = &RFC3164{}
	p.rfc.compileMatcher()
}

func (s *RFC3164) compileMatcher() {
	pri := `<([0-9]{1,3})>`
	ts := `([A-Za-z]+\s\d+(\s\d+)?\s\d+:\d+:\d+:\s)` // with or without year
	host := `([^ ]+)`
	msg := `(.+$)`
	s.matcher = regexp.MustCompile(pri + ts + `?` + host + `:\s` + msg)
}

func (s *RFC3164) parse(raw []byte, result *map[string]interface{}) {
	m := s.matcher.FindStringSubmatch(string(raw))
	if m == nil || len(m) < 4 {
		return
	}
	pri, _ := strconv.Atoi(m[1])
	*result = map[string]interface{}{
		"priority": pri,
	}
	if len(m) == 5 {
		(*result)["timestamp"] = m[2]
		(*result)["host"] = m[3]
		(*result)["message"] = m[4]
	} else {
		(*result)["host"] = m[3]
		(*result)["message"] = m[4]
	}
}

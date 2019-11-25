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
	"fmt"
	"strings"
)

var (
	fmtsByStandard    = []string{"rfc5424", "rfc3164"}
	fmtsByStandardNum = []string{"5424", "3164"}
)

// ValidFormat returns if the given format matches one of the possible formats.
func ValidFormat(format string) bool {
	for _, f := range append(fmtsByStandard, fmtsByStandardNum...) {
		if f == format {
			return true
		}
	}
	return false
}

// A Parser parses the raw input as a map with a timestamp field.
type Parser struct {
	fmt       string
	Raw       []byte
	Result    map[string]interface{}
	rfc       RFC
	delimiter Delimiter
}

// NewParser returns a new Parser instance.
func NewParser(f string) (*Parser, error) {
	if !ValidFormat(f) {
		return nil, fmt.Errorf("%s is not a valid format", f)
	}

	p := &Parser{}
	p.detectFmt(strings.TrimSpace(strings.ToLower(f)))
	switch p.fmt {
	case "rfc5424":
		p.newRFC5424Parser()
		p.delimiter = NewRFC5424Delimiter(msgBufSize)
		break
	case "rfc3164":
		p.newRFC3164Parser()
		p.delimiter = NewRFC3164Delimiter(msgBufSize)
		break
	}
	return p, nil
}

// Reads the given format and detects its internal name.
func (p *Parser) detectFmt(f string) {
	for i, v := range fmtsByStandardNum {
		if f == v {
			p.fmt = fmtsByStandard[i]
			return
		}
	}
	for _, v := range fmtsByStandard {
		if f == v {
			p.fmt = v
			return
		}
	}
	p.fmt = fmtsByStandard[0]
	return
}

// Parse the given byte slice.
func (p *Parser) Parse(b []byte) bool {
	p.Result = map[string]interface{}{}
	p.Raw = b
	p.rfc.parse(p.Raw, &p.Result)
	if len(p.Result) == 0 {
		return false
	}
	return true
}

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

// Unknown represents a parser for Unknown-compliant log messages
type Unknown struct {
	fortinet *regexp.Regexp
}

func (p *Parser) newUnknownParser() {
	p.rfc = &Unknown{}
	p.rfc.compileMatcher()
}

func (s *Unknown) compileMatcher() {
	// common
	pri := `<(-?[0-9]{1,3})>`
	s.compileFortinet(pri)
}

func (s *Unknown) compileFortinet(pri string) {
	// fortinet
	// date := `date=([^ ]+)`
	// time := `time=([^ ]+)`
	// devName := `devname=([^ ]+)`
	// devID := `devid=([^ ]+)`
	// logID := `logid=([^ ]+)`
	// t := `type=([^ ]+)`
	// subtype := `subtype=([^ ]+)`
	// level := `level=([^ ]+)`
	// vd := `vd=([^ ]+)`
	// logDesc := `logdesc=\"(.+?)\"`
	// action := `action=\"(.+?)\"`
	// cpu := `cpu=([0-9]+)`
	// mem := `mem=([0-9]+)`
	// totalSession := `totalsession=([0-9]+)`
	// disk := `disk=([0-9]+)`
	// bandwidth := `bandwidth=([^ ]+)`
	// ui := `ui=([^ ]+)`
	// msg := `msg=\"(.+?)\"`
	// eg. <185>date=2017-04-08 time=15:35:01 devname=fortinet239 devid=FGVMEV0000000000 logid=0100032400 type=event subtype=system level=alert vd=root logdesc="Configuration changed" user="admin" ui=ssh(192.168.1.160) msg="Configuration is changed in the admin session"
	s.fortinet = regexp.MustCompile(`(?m)<(-?[0-9]{1,3})>((([^ ]+)=([^ \"]+)\s)|(([^ ]+)=\"(.+?)\"\s?))+`)

}

func (s *Unknown) parse(raw []byte, result *map[string]interface{}) {
	// match fortinet
	str := string(raw)
	m := s.fortinet.FindStringSubmatch(string(raw))
	if m != nil {
		// do
		pri, _ := strconv.Atoi(m[1])
		*result = map[string]interface{}{
			"priority": pri,
		}
		devnameMatch := regexp.MustCompile(`devname=([^ ]+)`).FindStringSubmatch(str)
		// (*result)["timestamp"] = m[2] + " " + m[3]
		(*result)["hostname"] = devnameMatch[1]
		(*result)["host"] = devnameMatch[1]
		(*result)["vendor"] = "fortinet"
		return
	}
	// other match
}

package syslog

import (
	"fmt"
	"regexp"
	"strconv"
)

// RFC3164 represents a parser for RFC3164-compliant log messages
// BUT made some modifications to make it compatible with juniper syslog messages
type RFC3164 struct {
	matcher []*regexp.Regexp
}

func (p *Parser) newRFC3164Parser() {
	p.rfc = &RFC3164{}
	p.rfc.compileMatcher()
}

func (s *RFC3164) compileMatcher() {
	s.matcher = make([]*regexp.Regexp, 7)
	pri := `<([0-9]{1,3})>`
	ts := `([A-Za-z]+\s\d+(\s\d+)?\s\d+:\d+:\d+)`  // without year
	chost := `([:alnum:]]{1,}(-[[:alnum:]]+){0,})` // cisco hostname
	jhost := `([[:alnum:]._-]+)`                   // juniper hostname
	// uuid := `([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})`
	uuid := `([^ =:]+)`
	capp := `(%\w+-\d+-\d+)` // %ASA-5-611103
	japp := `(\w+\[\d+\])`   // sshd[8144]
	msg := `(.+$)`
	// cisco host only
	s.matcher[0] = regexp.MustCompile(pri + ts + `\s` + chost + `\s:\s` + capp + `:\s` + msg)
	// cisco uuid only
	s.matcher[1] = regexp.MustCompile(pri + ts + `\s` + uuid + `\s:\s` + capp + `:\s` + msg)
	// juniper host+uuid
	s.matcher[2] = regexp.MustCompile(pri + ts + `\s` + jhost + `\s` + uuid + `:\s` + japp + `:\s` + msg)
	// juniper host only
	s.matcher[3] = regexp.MustCompile(pri + ts + `\s` + jhost + `:\s` + japp + `:\s` + msg)

	// fortinet
	leading := `(?s)`
	// pri := `<([0-9]{1,3})>`
	// date := `date=([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))`
	date := `date=([0-9-]+)`
	// time := `time=(2[0-3]|[01]?[0-9]):([0-5]?[0-9]):([0-5]?[0-9])`
	time := `time=([0-9:]{8})`
	devName := `devname=([0-9A-Za-z][0-9A-Za-z_-]*)`
	devID := `devid=([0-9A-Za-z][0-9A-Za-z_-]*)`
	logID := `logid=([0-9]+)`
	logType := `type=([a-z]+)`
	subType := `subtype=([a-z]+)`
	level := `level=(emergency|alert|critical|error|warning|notification|information|debug)`
	vd := `vd=([a-zA-Z]+)`
	logDesc := `logdesc="([[:alnum:]\s]+)"`
	// action := `action=([[:alpha:]]+)`
	// user := `user="([[:alnum:]]+)"`
	// ui := `ui=ssh\(([[:digit:]]+(\.[[:digit:]]+){3})\)` //
	// quoteMsg := `msg="(.*)"$`
	trailing := `(.*)$`
	mstr := leading +
		pri +
		date + `\s` +
		time + `\s` +
		devName + `\s` +
		devID + `\s` +
		logID + `\s` +
		logType + `\s` +
		subType + `\s` +
		level + `\s` +
		vd + `\s` +
		logDesc + `\s` +
		// user + `\s` +
		// ui + `\s` +
		// quoteMsg
		trailing
	s.matcher[4] = regexp.MustCompile(mstr)

	nexusPattern := pri +
		uuid + `:\s` +
		`(\d+\s[A-Za-z]+\s+\d\s\d+:\d+:\d+\s[\w]+)` + `:\s` +
		trailing
	s.matcher[5] = regexp.MustCompile(nexusPattern)

	iosPattern := pri +
		`(\d+)` + `:\s` +
		`s([^ =:]+):\` + `:\s\*` +
		`(\w+\s\s\d\s\d+:\d+:\d+\.\d+)` + `:\s` +
		trailing
	s.matcher[6] = regexp.MustCompile(iosPattern)

	paloaltoPattern := pri + ts + `\s` +
		`([a-zA-Z0-9\._-]+)` + `\s` + // hostname
		`([0-9]+)` + `,` +
		`([0-9\/])+\s\d+:\d+:\d+` + `,` + // with year
		trailing

	s.matcher[7] = regexp.MustCompile(paloaltoPattern)

}

func (s *RFC3164) parse(raw []byte, result *map[string]interface{}) {
	// fmt.Println(string(raw))
	for i, v := range s.matcher[0:4] {
		m := v.FindStringSubmatch(string(raw))
		if len(m) == 0 {
			continue
		}
		pri, _ := strconv.Atoi(m[1])

		*result = map[string]interface{}{
			"priority":  pri,
			"timestamp": m[2],
		}
		if i == 0 || i == 1 || i == 3 {
			(*result)["identifier"] = m[4]
			(*result)["app"] = m[5]
			(*result)["message"] = m[6]
		} else if i == 2 {
			(*result)["identifier"] = m[5]
			(*result)["app"] = m[6]
			(*result)["message"] = m[7]
		}
		return
	}
	// fortinet
	m := s.matcher[4].FindStringSubmatch(string(raw))
	for i, v := range m {
		fmt.Println(i, ":", v)
	}
	if len(m) != 0 {
		pri, _ := strconv.Atoi(m[1])
		*result = map[string]interface{}{
			"priority":   pri,
			"timestamp":  m[2] + " " + m[3],
			"identifier": m[4],
			"message":    m[11],
		}
		return
	}

	n := s.matcher[5].FindStringSubmatch(string(raw))
	for i, v := range n {
		fmt.Println(i, ":", v)
	}
	if len(m) != 0 {
		pri, _ := strconv.Atoi(n[1])
		*result = map[string]interface{}{
			"priority":   pri,
			"timestamp":  n[3],
			"identifier": n[2],
			"message":    n[4],
		}
		return
	}

	o := s.matcher[6].FindStringSubmatch(string(raw))
	for i, v := range o {
		fmt.Println(i, ":", v)
	}
	if len(o) != 0 {
		pri, _ := strconv.Atoi(o[1])
		*result = map[string]interface{}{
			"priority":   pri,
			"timestamp":  n[4],
			"identifier": n[3],
			"message":    n[5],
		}
	}

	p := s.matcher[7].FindStringSubmatch(string(raw))
	for i, v := range p {
		fmt.Println(i, ":", v)
	}
	if len(p) != 0 {
		pri, _ := strconv.Atoi(p[1])
		*result = map[string]interface{}{
			"priority":   pri,
			"timestamp":  p[5],
			"identifier": p[3],
			"message":    p[6],
		}
	}
}

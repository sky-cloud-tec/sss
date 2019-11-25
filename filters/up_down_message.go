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

package filters

import (
	"regexp"

	"github.com/sky-cloud-tec/sss/common"
)

// UpDownMessage filter
type UpDownMessage struct {
}

// NewUpDownMessage return up down message filter
func NewUpDownMessage() Filter {
	return &UpDownMessage{}
}

// Match return true if message matches the condition
func (c *UpDownMessage) Match(msg *common.Message) bool {
	nexusUpDownPattern := `Interface ([\w\d\/]+) is (up|down)`
	message := msg.Parsed["message"].(string)
	m := regexp.MustCompile(nexusUpDownPattern).FindStringSubmatch(message)
	if len(m) > 0 {
		msg.Parsed["interface"] = m[1]
		msg.Parsed["state"] = m[2]
		return true
	}
	iosUpDownPattern := `Line protocol on Interface ([\w\d\/]+), changed state to (up|down)`
	n := regexp.MustCompile(iosUpDownPattern).FindStringSubmatch(message)
	if len(n) > 0 {
		msg.Parsed["interface"] = m[1]
		msg.Parsed["state"] = m[2]
		return true
	}
	return false
}

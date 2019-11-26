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
	"strings"

	"github.com/sky-cloud-tec/sss/common"
)

// NewCommitMessage return commit message filter
func NewCommitMessage() Filter {
	return &CommitMessage{}

}

// CommitMessage filter
type CommitMessage struct {
}

// Match return true if message matches the condition
func (c *CommitMessage) Match(msg *common.Message) bool {
	message := msg.Parsed["message"].(string)
	return strings.Contains(message, "'write memory' command") || // cisco asa
		strings.Contains(message, "commit complete") || // juniper srx
		strings.Contains(strings.ToLower(message), "attribute configured") || // fortinet
		strings.Contains(message, "Configuration is changed") || // fortinet
		strings.Contains(strings.ToLower(message), "configured from") || // NX-OS || IOS
		(strings.Contains(message, "commit") && strings.Contains(message, "Submitted"))
}

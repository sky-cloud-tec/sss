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
		strings.Contains(strings.ToLower(message), "configured from") || // NX-OS || IOS
		(strings.Contains(message, "commit") && strings.Contains(message, "Submitted"))
}

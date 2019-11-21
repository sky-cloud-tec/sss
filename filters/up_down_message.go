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

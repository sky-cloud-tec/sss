package filters

import (
	"github.com/sky-cloud-tec/sss/common"
)

// Filter interface
type Filter interface {
	Match(msg *common.Message) bool
}

// FilterMap map filter name to filter instance
var FilterMap map[string]Filter

func init() {
	FilterMap = make(map[string]Filter, 0)
	FilterMap["commit"] = NewCommitMessage()
	FilterMap["up_down"] = NewUpDownMessage()
}

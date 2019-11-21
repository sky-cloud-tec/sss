package consumers

import (
	"github.com/sky-cloud-tec/sss/common"
)

// Consumer consume messages
type Consumer interface {
	C() chan *common.Message
	Consume()
}

package consumers

import (
	"time"

	"github.com/grafana/loki/pkg/promtail/client"
	lokiflag "github.com/grafana/loki/pkg/util/flagext"
	"github.com/prometheus/common/config"
	"github.com/sky-cloud-tec/sss/common"
	"github.com/songtianyi/rrframework/logs"
	"github.com/prometheus/common/model"
)

type lokiConsumer struct {
	c      chan *common.Message
	client client.Client
}

// NewLokiConsumer create loki consumer
func NewLokiConsumer(url string) (Consumer, error) {
	cfg := client.Config{
		URL:            url,
		BatchWait:      100 * time.Millisecond,
		BatchSize:      10,
		Client:         config.HTTPClientConfig{},
		BackoffConfig:  util.BackoffConfig{MinBackoff: 1 * time.Millisecond, MaxBackoff: 2 * time.Millisecond, MaxRetries: 2},
		ExternalLabels: lokiflag.LabelSet{},
		Timeout:        1 * time.Second,
	}
	logs.Info("creating loki consumer...")
	// Obtain a client. You can also provide your own HTTP client here.
	c, err := client.New(cfg, logs.GetLogger())
	return &lokiConsumer{
		client: c,
		c:      make(chan *common.Message, 0),
	}, nil
}

func (e *lokiConsumer) C() chan *common.Message {
	return e.c

}

func (e *lokiConsumer) Consume() {
	for {
		select {
		case msg := <-e.c:
			if err := e.do(msg); err != nil {
				logs.Error(err)
			}
		}
	}
}

func (e *lokiConsumer) do(msg *common.Message) error {
	err = return e.client.Handle(model.LabelSet{}, msg.ReceptionTime, msg.Text)
	return nil
}

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

package consumers

// import (
// 	"time"

// 	"github.com/grafana/loki/pkg/promtail/client"
// 	lokiflag "github.com/grafana/loki/pkg/util/flagext"
// 	"github.com/prometheus/common/config"
// 	"github.com/sky-cloud-tec/sss/common"
// 	"github.com/songtianyi/rrframework/logs"
// 	"github.com/prometheus/common/model"
// )

// type lokiConsumer struct {
// 	c      chan *common.Message
// 	client client.Client
// }

// // NewLokiConsumer create loki consumer
// func NewLokiConsumer(url string) (Consumer, error) {
// 	cfg := client.Config{
// 		URL:            url,
// 		BatchWait:      100 * time.Millisecond,
// 		BatchSize:      10,
// 		Client:         config.HTTPClientConfig{},
// 		BackoffConfig:  util.BackoffConfig{MinBackoff: 1 * time.Millisecond, MaxBackoff: 2 * time.Millisecond, MaxRetries: 2},
// 		ExternalLabels: lokiflag.LabelSet{},
// 		Timeout:        1 * time.Second,
// 	}
// 	logs.Info("creating loki consumer...")
// 	// Obtain a client. You can also provide your own HTTP client here.
// 	c, err := client.New(cfg, logs.GetLogger())
// 	return &lokiConsumer{
// 		client: c,
// 		c:      make(chan *common.Message, 0),
// 	}, nil
// }

// func (e *lokiConsumer) C() chan *common.Message {
// 	return e.c

// }

// func (e *lokiConsumer) Consume() {
// 	for {
// 		select {
// 		case msg := <-e.c:
// 			if err := e.do(msg); err != nil {
// 				logs.Error(err)
// 			}
// 		}
// 	}
// }

// func (e *lokiConsumer) do(msg *common.Message) error {
// 	err = return e.client.Handle(model.LabelSet{}, msg.ReceptionTime, msg.Text)
// 	return nil
// }

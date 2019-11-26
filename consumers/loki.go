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

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sky-cloud-tec/sss/common"
	"github.com/songtianyi/rrframework/logs"
)

type lokiConsumer struct {
	c   chan *common.Message
	url string
}

const (
	// APIPush push url
	APIPush = "/loki/api/v1/push"
)

// PushRequest POST /loki/api/v1/push body struct
type PushRequest struct {
	Streams []Stream `json:"streams"`
}

// Stream struct
type Stream struct {
	Labels  map[string]string `json:"stream"`
	Entries [][]string        `json:"values"`
}

// NewLokiConsumer create loki consumer
func NewLokiConsumer(url string) (Consumer, error) {
	logs.Info("creating loki consumer...")
	return &lokiConsumer{
		url: url,
		c:   make(chan *common.Message, 0),
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
	pr := &PushRequest{
		Streams: make([]Stream, 0),
	}
	stream := Stream{
		Labels:  map[string]string{"type": "firewall", "source": msg.SourceIP},
		Entries: [][]string{{strconv.FormatInt(msg.ReceptionTime.UnixNano(), 10), msg.Text}},
	}
	pr.Streams = append(pr.Streams, stream)
	b, _ := json.Marshal(pr)
	client := &http.Client{}
	req, err := http.NewRequest("POST", e.url+APIPush, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logs.Debug(string(body))
	return nil
}

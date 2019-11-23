package consumers

import (
	"context"

	"github.com/olivere/elastic"
	"github.com/sky-cloud-tec/sss/common"
	"github.com/songtianyi/rrframework/logs"
)

type esConsumer struct {
	client *elastic.Client
	c      chan *common.Message
}

// NewESConsumer create ealsticsearch consumer
func NewESConsumer(url, username, password string) (Consumer, error) {
	logs.Info("creating es consumer...")
	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetBasicAuth(username, password))
	if err != nil {
		// Handle error
		return nil, err
	}
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(url).Do(context.Background())
	if err != nil {
		// Handle error
		return nil, err
	}
	logs.Info("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	return &esConsumer{
		client: client,
		c:      make(chan *common.Message, 0),
	}, nil
}

func (e *esConsumer) C() chan *common.Message {
	return e.c

}

func (e *esConsumer) Consume() {
	for {
		select {
		case msg := <-e.c:
			if err := e.do(msg); err != nil {
				logs.Error(err)
			}
		}
	}
}

func (e *esConsumer) do(msg *common.Message) error {
	put2, err := e.client.Index().
		Index("firewall_syslog").
		Type("syslog").
		// BodyString(msg.Text).
		BodyJson(msg.Parsed).
		Do(context.Background())
	if err != nil {
		return err
	}
	logs.Debug(put2)
	return nil
}

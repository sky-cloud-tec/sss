package syslog

import (
	"fmt"

	"github.com/sky-cloud-tec/sss/common"
	"github.com/sky-cloud-tec/sss/consumers"
	"github.com/sky-cloud-tec/sss/filters"
)

// Pipe struct
type Pipe struct {
	c         chan *common.Message
	filters   []filters.Filter
	consumers []consumers.Consumer
}

// NewPipe create a pipe and return
func NewPipe() *Pipe {
	return &Pipe{c: make(chan *common.Message, 0), filters: make([]filters.Filter, 0)}
}

// C return pipe channel
func (p *Pipe) C() chan *common.Message {
	return p.c
}

// Apply a filter
func (p *Pipe) Apply(f filters.Filter) *Pipe {
	p.filters = append(p.filters, f)
	return p
}

// Register consumers
func (p *Pipe) Register(c consumers.Consumer) *Pipe {
	p.consumers = append(p.consumers, c)
	return p
}

// Open pipe
func (p *Pipe) Open() {
	for {
		select {
		case msg := <-p.c:
			fmt.Println(msg)
			if len(p.filters) == 0 {
				// no filter applied send all msgs to consumers
				for _, consumer := range p.consumers {
					consumer.C() <- msg
				}
				break
			}
			for _, filter := range p.filters {
				// filter | filter
				if filter.Match(msg) {
					// send to consumers
					for _, consumer := range p.consumers {
						consumer.C() <- msg
					}
					break
				}
			}
		}
	}
}

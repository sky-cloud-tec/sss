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

package syslog

import (
	"github.com/sky-cloud-tec/sss/common"
	"github.com/sky-cloud-tec/sss/consumers"
	"github.com/sky-cloud-tec/sss/filters"
	"github.com/songtianyi/rrframework/logs"
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
			if len(p.filters) == 0 {
				// no filter applied send all msgs to consumers
				logs.Debug("[accept]", msg)
				for _, consumer := range p.consumers {
					consumer.C() <- msg
				}
				break
			}
			for _, filter := range p.filters {
				// filter | filter
				if filter.Match(msg) {
					logs.Debug("[accept]", msg)
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

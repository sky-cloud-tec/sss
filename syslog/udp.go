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
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sky-cloud-tec/sss/common"
	"github.com/songtianyi/rrframework/logs"
)

// UDPCollector represents a network collector that accepts UDP packets.
type UDPCollector struct {
	format string
	addr   *net.UDPAddr
}

// Start instructs the UDPCollector to start reading packets from the interface.
func (s *UDPCollector) Start(c chan<- *common.Message) error {
	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}

	parser, err := NewParser(s.format)
	if err != nil {
		panic(fmt.Sprintf("failed to create UDP parser:%s", err.Error()))
	}

	go func() {
		buf := make([]byte, msgBufSize)
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			log := strings.Trim(string(buf[:n]), "\r\n")
			logs.Debug("[raw]", log)
			if parser.Parse(bytes.NewBufferString(log).Bytes()) {
				c <- &common.Message{
					Text:          log,
					Parsed:        parser.Result,
					ReceptionTime: time.Now().UTC(),
					SourceIP:      addr.String(),
				}
			} else {
				logs.Error("parse raw msg", log, "error")
				panic(err)
			}
		}
	}()
	return nil
}

// Addr returns the net.Addr to which the UDP collector is bound.
func (s *UDPCollector) Addr() net.Addr {
	return s.addr
}

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
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sky-cloud-tec/sss/common"
	"github.com/songtianyi/rrframework/logs"
)

const (
	newlineTimeout = time.Duration(1000 * time.Millisecond)
	msgBufSize     = 256
)

// TCPCollector represents a network collector that accepts and handler TCP connections.
type TCPCollector struct {
	iface  string
	format string

	addr      net.Addr
	tlsConfig *tls.Config
}

// Start instructs the TCPCollector to bind to the interface and accept connections.
func (s *TCPCollector) Start(c chan<- *common.Message) error {
	var ln net.Listener
	var err error
	if s.tlsConfig == nil {
		ln, err = net.Listen("tcp", s.iface)
	} else {
		ln, err = tls.Listen("tcp", s.iface, s.tlsConfig)
	}
	if err != nil {
		return err
	}
	s.addr = ln.Addr()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn, c)
		}
	}()
	return nil
}

func (s *TCPCollector) handleConnection(conn net.Conn, c chan<- *common.Message) {
	defer func() {
		conn.Close()
	}()

	parser, err := NewParser(s.format)
	if err != nil {
		panic(fmt.Sprintf("failed to create TCP connection parser:%s", err.Error()))
	}

	reader := bufio.NewReader(conn)
	var log string
	var match bool

	for {
		conn.SetReadDeadline(time.Now().Add(newlineTimeout))
		b, err := reader.ReadByte()
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				logs.Notice("tcpConnReadTimeout")
			} else if err == io.EOF {
				logs.Error("tcpConnReadEOF")
			} else {
				logs.Error("tcpConnUnrecoverError")
				return
			}
			log, match = parser.delimiter.Vestige()
		} else {
			log, match = parser.delimiter.Push(b)
		}

		// Log line available?
		if match {
			if parser.Parse(bytes.NewBufferString(log).Bytes()) {
				c <- &common.Message{
					Text:          string(parser.Raw),
					Parsed:        parser.Result,
					ReceptionTime: time.Now().UTC(),
					SourceIP:      conn.RemoteAddr().String(),
				}
			} else {
				// Zero tolerance :)
				panic(err)
			}
		}

		// Was the connection closed?
		if err == io.EOF {
			return
		}
	}
}

// Addr returns the net.Addr that the Collector is bound to, in a race-say manner.
func (s *TCPCollector) Addr() net.Addr {
	return s.addr
}

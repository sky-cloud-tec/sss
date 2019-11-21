package syslog

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/sky-cloud-tec/sss/common"
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
			// fmt.Println(log)
			if parser.Parse(bytes.NewBufferString(log).Bytes()) {
				c <- &common.Message{
					Text:          log,
					Parsed:        parser.Result,
					ReceptionTime: time.Now().UTC(),
					SourceIP:      addr.String(),
				}
			}
		}
	}()
	return nil
}

// Addr returns the net.Addr to which the UDP collector is bound.
func (s *UDPCollector) Addr() net.Addr {
	return s.addr
}

package syslog

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/sky-cloud-tec/sss/common"
)

// Collector specifies the interface all network collectors must implement.
type Collector interface {
	Start(chan<- *common.Message) error
	Addr() net.Addr
}

// NewCollector returns a network collector of the specified type, that will bind
// to the given inteface on Start(). If config is non-nil, a secure Collector will
// be returned. Secure Collectors require the protocol be TCP.
func NewCollector(proto, iface, format string, tlsConfig *tls.Config) (Collector, error) {
	if strings.ToLower(proto) == "tcp" {
		return &TCPCollector{
			iface:     iface,
			format:    format,
			tlsConfig: tlsConfig,
		}, nil
	} else if strings.ToLower(proto) == "udp" {
		addr, err := net.ResolveUDPAddr("udp", iface)
		if err != nil {
			return nil, err
		}
		return &UDPCollector{addr: addr, format: format}, nil
	}
	return nil, fmt.Errorf("unsupport collector protocol")
}

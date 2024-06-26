package outbound

import (
	"encoding/json"
	"net"
	"syscall"

	N "github.com/carlos19960601/ClashV/common/net"
	"github.com/carlos19960601/ClashV/common/utils"
	"github.com/carlos19960601/ClashV/component/dialer"
	C "github.com/carlos19960601/ClashV/constant"
)

type Base struct {
	name   string
	addr   string
	tp     C.AdapterType
	udp    bool
	id     string
	prefer C.DNSPrefer
}

type BasicOption struct {
	Name string
}

// Id implements C.ProxyAdapter
func (b *Base) Id() string {
	if b.id == "" {
		b.id = utils.NewUUIDV6().String()
	}

	return b.id
}

// Name implements C.ProxyAdapter
func (b *Base) Name() string {
	return b.name
}

// Type implements C.ProxyAdapter
func (b *Base) Type() C.AdapterType {
	return b.tp
}

// Addr implements C.ProxyAdapter
func (b *Base) Addr() string {
	return b.addr
}

// SupportUDP implements C.ProxyAdapter
func (b *Base) SupportUDP() bool {
	return b.udp
}

func (b *Base) DialOptions(opts ...dialer.Option) []dialer.Option {
	return opts
}

// MarshalJSON implements C.ProxyAdapter
func (b *Base) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type": b.Type().String(),
		"id":   b.Id(),
	})
}

func NewConn(c net.Conn, a C.ProxyAdapter) C.Conn {
	if _, ok := c.(syscall.Conn); !ok {
		c = N.NewDeadlineConn(c)
	}
	return &conn{N.NewExtendedConn(c), []string{a.Name(), parseRemoteDestination(a.Addr())}}
}

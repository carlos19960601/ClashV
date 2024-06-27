package outbound

import (
	"encoding/json"
	"net"
	"strings"
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
	Interface string `proxy:"interface-name,omitempty" group:"interface-name,omitempty"`
}

type BaseOption struct {
	Name      string
	Addr      string
	Type      C.AdapterType
	UDP       bool
	Interface string
	Perfer    C.DNSPrefer
}

func NewBase(opt BaseOption) *Base {
	return &Base{
		name: opt.Name,
		addr: opt.Addr,
		tp:   opt.Type,
	}
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

func (b *Base) Unwrap(metadata *C.Metadata, touch bool) C.Proxy {
	return nil
}

// MarshalJSON implements C.ProxyAdapter
func (b *Base) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type": b.Type().String(),
		"id":   b.Id(),
	})
}

type conn struct {
	N.ExtendedConn
	chain                   C.Chain
	actualRemoteDestination string
}

func NewConn(c net.Conn, a C.ProxyAdapter) C.Conn {
	if _, ok := c.(syscall.Conn); !ok { // exclusion system conn like *net.TCPConn
		c = N.NewDeadlineConn(c) // most conn from outbound can't handle readDeadline correctly
	}
	return &conn{N.NewExtendedConn(c), []string{a.Name()}, parseRemoteDestination(a.Addr())}
}

// Chains implements C.Connection
func (c *conn) Chains() C.Chain {
	return c.chain
}

func parseRemoteDestination(addr string) string {
	if dst, _, err := net.SplitHostPort(addr); err == nil {
		return dst
	} else {
		if addrError, ok := err.(*net.AddrError); ok && strings.Contains(addrError.Err, "missing port") {
			return dst
		} else {
			return ""
		}
	}
}

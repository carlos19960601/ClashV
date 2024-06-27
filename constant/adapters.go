package constant

import (
	"context"
	"fmt"

	N "github.com/carlos19960601/ClashV/common/net"
	"github.com/carlos19960601/ClashV/component/dialer"
)

// Adapter Type
const (
	Direct AdapterType = iota
	Reject
	RejectDrop
	Compatible
	Pass
	Dns

	Relay
	Selector
	Fallback
	URLTest
	LoadBalance

	Shadowsocks
	ShadowsocksR
	Snell
	Socks5
	Http
	Vmess
	Vless
	Trojan
	Hysteria
	Hysteria2
	WireGuard
	Tuic
	Ssh
)

type AdapterType int

func (at AdapterType) String() string {
	switch at {
	case Direct:
		return "Direct"
	case Reject:
		return "Reject"
	case RejectDrop:
		return "RejectDrop"
	case Compatible:
		return "Compatible"
	case Pass:
		return "Pass"
	case Dns:
		return "Dns"
	case Shadowsocks:
		return "Shadowsocks"
	case ShadowsocksR:
		return "ShadowsocksR"
	case Snell:
		return "Snell"
	case Socks5:
		return "Socks5"
	case Http:
		return "Http"
	case Vmess:
		return "Vmess"
	case Vless:
		return "Vless"
	case Trojan:
		return "Trojan"
	case Hysteria:
		return "Hysteria"
	case Hysteria2:
		return "Hysteria2"
	case WireGuard:
		return "WireGuard"
	case Tuic:
		return "Tuic"

	case Relay:
		return "Relay"
	case Selector:
		return "Selector"
	case Fallback:
		return "Fallback"
	case URLTest:
		return "URLTest"
	case LoadBalance:
		return "LoadBalance"
	case Ssh:
		return "Ssh"
	default:
		return "Unknown"
	}
}

type Proxy interface {
	ProxyAdapter
}

type ProxyAdapter interface {
	Name() string
	Type() AdapterType
	Addr() string
	SupportUDP() bool
	MarshalJSON() ([]byte, error)

	DialContext(ctx context.Context, metadata *Metadata, opts ...dialer.Option) (Conn, error)

	// Unwrap 从代理组中获取代理
	Unwrap(metadata *Metadata, touch bool) Proxy
}

type Conn interface {
	N.ExtendedConn
	Connection
}

const (
	DefaultTCPTimeout = dialer.DefaultTCPTimeout
	DefaultUDPTimeout = dialer.DefaultUDPTimeout
	DefaultDropTime   = 12 * DefaultTCPTimeout
	DefaultTLSTimeout = DefaultTCPTimeout
	DefaultTestURL    = "https://www.gstatic.com/generate_204"
)

type Chain []string

func (c Chain) String() string {
	switch len(c) {
	case 0:
		return ""
	case 1:
		return c[0]
	default:
		return fmt.Sprintf("%s[%s]", c[len(c)-1], c[0])
	}
}

type Connection interface {
	Chains() Chain
}

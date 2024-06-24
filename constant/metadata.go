package constant

import (
	"net"
	"net/netip"
	"strconv"
)

const (
	HTTP Type = iota
	HTTPS
	SOCKS4
	SOCKS5
	SHADOWSOCKS
	VMESS
	REDIR
	TPROXY
	TUNNEL
	TUN
	TUIC
	HYSTERIA2
	INNER
)

type Type int

func (t Type) String() string {
	switch t {
	case HTTP:
		return "HTTP"
	case HTTPS:
		return "HTTPS"
	case SOCKS4:
		return "Socks4"
	case SOCKS5:
		return "Socks5"
	case SHADOWSOCKS:
		return "ShadowSocks"
	case VMESS:
		return "Vmess"
	case REDIR:
		return "Redir"
	case TPROXY:
		return "TProxy"
	case TUNNEL:
		return "Tunnel"
	case TUN:
		return "Tun"
	case TUIC:
		return "Tuic"
	case HYSTERIA2:
		return "Hysteria2"
	case INNER:
		return "Inner"
	default:
		return "Unknown"
	}
}

// Socks addr type
const (
	TCP NetWork = iota
	UDP
	ALLNet
	InvalidNet = 0xff
)

type NetWork int

func (n NetWork) String() string {
	switch n {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	case ALLNet:
		return "all"
	default:
		return "invalid"
	}
}

type Metadata struct {
	NetWork NetWork    `json:"network"`
	Type    Type       `json:"type"`
	SrcIP   netip.Addr `json:"sourceIP"`
	DstIP   netip.Addr `json:"destinationIP"`
	SrcPort uint16     `json:"sourcePort,string"`
	DstPort uint16     `json:"destinationPort,string"`
	Host    string     `json:"host"`
	InIP    netip.Addr `json:"inboundIP"`
	InPort  uint16     `json:"inboundPort,string"`
}

func (m *Metadata) SetRemoteAddr(addr net.Addr) error {
	if addr == nil {
		return nil
	}

	return m.SetRemoteAddress(addr.String())
}

func (m *Metadata) SetRemoteAddress(rawAddress string) error {
	host, port, err := net.SplitHostPort(rawAddress)
	if err != nil {
		return err
	}

	var uint16Port uint16
	if port, err := strconv.ParseUint(port, 10, 16); err == nil {
		uint16Port = uint16(port)
	}

	if ip, err := netip.ParseAddr(host); err != nil {
		m.Host = host
		m.DstIP = netip.Addr{}
	} else {
		m.Host = ""
		m.DstIP = ip.Unmap()
	}

	m.DstPort = uint16Port

	return nil
}

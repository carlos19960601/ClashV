package resolver

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"time"

	"github.com/carlos19960601/ClashV/component/trie"
)

var (
	DefaultResolver Resolver

	ProxyServerHostResolver Resolver

	DefaultHosts = NewHosts(trie.New[HostValue]())

	DefaultDNSTimeout = time.Second * 5
)

var (
	ErrIPNotFound = errors.New("couldn't find ip")
)

type Resolver interface {
	LookupIP(ctx context.Context, host string) (ips []netip.Addr, err error)
	LookupIPv4(ctx context.Context, host string) (ips []netip.Addr, err error)
	LookupIPv6(ctx context.Context, host string) (ips []netip.Addr, err error)
	Invalid() bool
}

func LookupIPv4ProxyServerHost(ctx context.Context, host string) ([]netip.Addr, error) {
	return nil, nil
}

func ResolveIP(ctx context.Context, host string) (netip.Addr, error) {
	return ResolveIPWithResolver(ctx, host, DefaultResolver)
}

func ResolveIPWithResolver(ctx context.Context, host string, r Resolver) (netip.Addr, error) {
	return netip.Addr{}, nil
}

func LookupIPWithResolver(ctx context.Context, host string, r Resolver) ([]netip.Addr, error) {
	if node, ok := DefaultHosts.Search(host, false); ok {
		return node.IPs, nil
	}

	if r != nil && r.Invalid() {

	}

	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if err != nil {
		return nil, err
	} else if len(ips) == 0 {
		return nil, ErrIPNotFound
	}

	return ips, nil
}

func LookupIPProxyServerHost(ctx context.Context, host string) ([]netip.Addr, error) {
	if ProxyServerHostResolver != nil {
		return LookupIPWithResolver(ctx, host, ProxyServerHostResolver)
	}

	return LookupIP(ctx, host)
}

func LookupIP(ctx context.Context, host string) ([]netip.Addr, error) {
	return LookupIPWithResolver(ctx, host, DefaultResolver)
}

func SortationAddr(ips []netip.Addr) (ipv4s, ipv6s []netip.Addr) {
	for _, v := range ips {
		if v.Unmap().Is4() {
			ipv4s = append(ipv4s, v)
		} else {
			ipv6s = append(ipv6s, v)
		}
	}

	return
}

package resolver

import (
	"context"
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

package resolver

import (
	"context"
	"net/netip"
)

var (
	DefaultResolver Resolver

	ProxyServerHostResolver Resolver
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

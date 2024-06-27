package dialer

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/carlos19960601/ClashV/component/resolver"
)

const (
	DefaultTCPTimeout = 5 * time.Second
	DefaultUDPTimeout = DefaultTCPTimeout
)

type dialFunc func(ctx context.Context, network string, ips []netip.Addr, port string, opt *option) (net.Conn, error)

var (
	actualSingleStackDialContext = serialSingleStackDialContext
	actualDualStackDialContext   = serialDualStackDialContext
)

func DialContext(ctx context.Context, network, address string, options ...Option) (net.Conn, error) {
	opt := applyOptions(options...)

	if opt.network == 4 || opt.network == 6 {
		if strings.Contains(network, "tcp") {
			network = "tcp"
		} else {
			network = "udp"
		}

		network = fmt.Sprintf("%s%d", network, opt.network)
	}

	ips, port, err := parseAddr(ctx, network, address, opt.resolver)
	if err != nil {
		return nil, err
	}

	switch network {
	case "tcp4", "tcp6":
		return actualSingleStackDialContext(ctx, network, ips, port, opt)
	case "tcp", "udp":
		return actualDualStackDialContext(ctx, network, ips, port, opt)
	default:
		return nil, ErrorInvalidedNetworkStack
	}
}

func parseAddr(ctx context.Context, network, address string, preferResolver resolver.Resolver) ([]netip.Addr, string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, "-1", err
	}

	var ips []netip.Addr
	switch network {
	case "tcp4", "udp4":
		if preferResolver == nil {
			ips, err = resolver.LookupIPv4ProxyServerHost(ctx, host)
		}
	default:
		if preferResolver == nil {
			ips, err = resolver.LookupIPProxyServerHost(ctx, host)
		} else {
			ips, err = resolver.LookupIPWithResolver(ctx, host, preferResolver)
		}
	}

	return ips, port, nil
}

func applyOptions(options ...Option) *option {
	opt := &option{}

	return opt
}

func serialSingleStackDialContext(ctx context.Context, network string, ips []netip.Addr, port string, opt *option) (net.Conn, error) {
	return serialDialContext(ctx, network, ips, port, opt)
}

func serialDualStackDialContext(ctx context.Context, network string, ips []netip.Addr, port string, opt *option) (net.Conn, error) {
	return dualStackDialContext(ctx, serialDialContext, network, ips, port, opt)
}

func serialDialContext(ctx context.Context, network string, ips []netip.Addr, port string, opt *option) (net.Conn, error) {

	return nil, nil
}

func dualStackDialContext(ctx context.Context, dialFn dialFunc, network string, ips []netip.Addr, port string, opt *option) (net.Conn, error) {
	ipv4s, ipv6s := resolver.SortationAddr(ips)
	if len(ipv4s) == 0 && len(ipv6s) == 0 {
		return nil, ErrorNoIpAddress
	}

	return nil, nil
}

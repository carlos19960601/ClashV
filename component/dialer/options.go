package dialer

import (
	"context"
	"net"

	"github.com/carlos19960601/ClashV/component/resolver"
)

type NetDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type option struct {
	network   int
	netDialer NetDialer
	resolver  resolver.Resolver
}

type Option func(opt *option)

func WithNetDialer(netDialer NetDialer) Option {
	return func(opt *option) {
		opt.netDialer = netDialer
	}
}

func WithOnlySingleStack(isIPv4 bool) Option {
	return func(opt *option) {
		if isIPv4 {
			opt.network = 4
		} else {
			opt.network = 6
		}
	}
}

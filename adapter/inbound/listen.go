package inbound

import (
	"context"
	"net"
)

var (
	lc = net.ListenConfig{}
)

func Listen(network, address string) (net.Listener, error) {
	return ListenContext(context.Background(), network, address)
}

func ListenContext(ctx context.Context, network, address string) (net.Listener, error) {
	return lc.Listen(ctx, network, address)
}

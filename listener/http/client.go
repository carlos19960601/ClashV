package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/carlos19960601/ClashV/adapter/inbound"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/log"
	"github.com/carlos19960601/ClashV/transport/socks5"
)

func newClient(srcConn net.Conn, tunnel C.Tunnel, additions ...inbound.Addition) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				log.Infoln("network: %s, address: %s", network, address)
				if network != "tcp" && network != "tcp4" && network != "tcp6" {
					return nil, errors.New("不支持的网络 " + network)
				}

				dstAddr := socks5.ParseAddr(address)
				if dstAddr == nil {
					return nil, socks5.ErrAddressNotSupported
				}

				return nil, nil
			},
		},
	}
}

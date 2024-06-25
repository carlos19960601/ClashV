package outbound

import (
	"context"
	"net"
	"strconv"

	"github.com/carlos19960601/ClashV/component/dialer"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/transport/trojan"
)

type Trojan struct {
	*Base
	option   *TrojanOption
	instance *trojan.Trojan
}

type TrojanOption struct {
	Name           string `proxy:"name"`
	Server         string `proxy:"name"`
	Port           int    `proxy:"port"`
	Password       string `proxy:"password"`
	UDP            bool   `proxy:"udp"`
	SNI            string `proxy:"sni"`
	SkipCertVerify bool   `proxy:"skip-cert-verify,omitempty"`
}

func NewTrojan(option TrojanOption) (C.Proxy, error) {
	addr := net.JoinHostPort(option.Server, strconv.Itoa(option.Port))

	tOption := &trojan.Option{
		Password:       option.Password,
		ServerName:     option.Server,
		SkipCertVerify: option.SkipCertVerify,
	}

	t := &Trojan{
		Base: &Base{
			name: option.Name,
			addr: addr,
		},
		instance: trojan.New(tOption),
		option:   &option,
	}

	return t, nil
}

// DialContext implements C.ProxyAdapter
func (t *Trojan) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.Conn, error) {
	return nil, nil
}

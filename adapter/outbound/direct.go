package outbound

import (
	"context"

	"github.com/carlos19960601/ClashV/component/dialer"
	C "github.com/carlos19960601/ClashV/constant"
)

type Direct struct {
	*Base
}

type DirectOption struct {
	BasicOption
	Name string `proxy:"name"`
}

func NewDirect() *Direct {
	return &Direct{
		Base: &Base{
			name: "DIRECT",
			tp:   C.Direct,
			udp:  true,
		},
	}
}

// DialContext implements C.ProxyAdapter
func (d *Direct) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.Conn, error) {
	return nil, nil
}

package outbound

import (
	"context"

	"github.com/carlos19960601/ClashV/component/dialer"
	"github.com/carlos19960601/ClashV/component/loopback"
	"github.com/carlos19960601/ClashV/component/resolver"
	C "github.com/carlos19960601/ClashV/constant"
)

type Direct struct {
	*Base
	loopBack *loopback.Detector
}

type DirectOption struct {
	BasicOption
	Name string `proxy:"name"`
}

func NewDirect() *Direct {
	return &Direct{
		Base: &Base{
			name:   "DIRECT",
			tp:     C.Direct,
			udp:    true,
			prefer: C.DualStack,
		},
		loopBack: loopback.NewDetector(),
	}
}

// DialContext implements C.ProxyAdapter
func (d *Direct) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.Conn, error) {
	if err := d.loopBack.CheckConn(metadata); err != nil {
		return nil, err
	}
	opts = append(opts, dialer.WithResolver(resolver.DefaultResolver))

	c, err := dialer.DialContext(ctx, "tcp", metadata.RemoteAddress(), d.Base.DialOptions(opts...)...)
	if err != nil {
		return nil, err
	}

	return d.loopBack.NewConn(NewConn(c, d)), nil
}

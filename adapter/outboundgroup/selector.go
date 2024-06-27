package outboundgroup

import (
	"context"

	"github.com/carlos19960601/ClashV/adapter/outbound"
	"github.com/carlos19960601/ClashV/component/dialer"
	C "github.com/carlos19960601/ClashV/constant"
)

type Selector struct {
	*GroupBase
	selected string
}

func NewSelector(option *GroupCommonOption) *Selector {
	return &Selector{
		GroupBase: NewGroupBase(GroupBaseOption{
			outbound.BaseOption{
				Name:      option.Name,
				Type:      C.Selector,
				Interface: option.Interface,
			},
		}),
	}
}

func (s *Selector) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.Conn, error) {
	return nil, errFormat
}

func (s *Selector) selectedProxy(touch bool) C.Proxy {
	proxies := s.GetProxies(touch)
	for _, proxy := range proxies {
		if proxy.Name() == s.selected {
			return proxy
		}
	}
	return proxies[0]
}

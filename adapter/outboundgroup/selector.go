package outboundgroup

import (
	"context"

	"github.com/carlos19960601/ClashV/adapter/outbound"
	"github.com/carlos19960601/ClashV/component/dialer"
	C "github.com/carlos19960601/ClashV/constant"
)

type Selector struct {
	*GroupBase
	selected bool
}

func NewSelector(option *GroupCommonOption) *Selector {
	return &Selector{
		GroupBase: NewGroupBase(GroupBaseOption{
			outbound.BasicOption{
				Name: option.Name,
			},
		}),
	}
}

func (s *Selector) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (C.Conn, error) {

}

func (s *Selector) selectedProxy(touch bool) C.Proxy {
	proxies := s.GetProxies(touch)
	for _, proxy := range proxies {
		
	}
}

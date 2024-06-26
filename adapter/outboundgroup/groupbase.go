package outboundgroup

import (
	"github.com/carlos19960601/ClashV/adapter/outbound"
	C "github.com/carlos19960601/ClashV/constant"
)

type GroupBase struct {
	*outbound.Base
	proxies [][]C.Proxy
}

type GroupBaseOption struct {
	outbound.BasicOption
}

func NewGroupBase(opt GroupBaseOption) *GroupBase {
	gb := &GroupBase{
		Base: outbound.NewBase(opt.BaseOption),
	}

	return gb
}

func (gb *GroupBase) GetProxies(touch bool) []C.Proxy {
	var proxies []C.Proxy

	for _, p := range gb.proxies {
		proxies = append(proxies, p...)
	}

	return proxies
}

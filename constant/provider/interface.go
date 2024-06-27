package provider

import "github.com/carlos19960601/ClashV/constant"

const (
	Proxy ProviderType = iota
	Rule
)

type ProviderType int

func (pt ProviderType) String() string {
	switch pt {
	case Proxy:
		return "Proxy"
	case Rule:
		return "Rule"
	default:
		return "Unknown"
	}
}

type Provider interface {
	Name() string
	Type() ProviderType
}

type ProxyProvider interface {
	Provider
	Proxies() []constant.Proxy
}

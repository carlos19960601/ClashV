package provider

import (
	"errors"

	C "github.com/carlos19960601/ClashV/constant"
	types "github.com/carlos19960601/ClashV/constant/provider"
)

const (
	ReservedName = "default"
)

func NewCompatibleProvider(name string, proxies []C.Proxy, hc *HealthCheck) (*CompatibleProvider, error) {
	if len(proxies) == 0 {
		return nil, errors.New("provider至少需要一个代理")
	}

	pd := &compatibleProvider{
		name:        name,
		proxies:     proxies,
		healthCheck: hc,
	}
	wrapper := &CompatibleProvider{pd}
	return wrapper, nil
}

type CompatibleProvider struct {
	*compatibleProvider
}

type compatibleProvider struct {
	name        string
	healthCheck *HealthCheck
	proxies     []C.Proxy
	version     uint32
}

func (cp *compatibleProvider) Name() string {
	return cp.name
}

func (cp *compatibleProvider) Type() types.ProviderType {
	return types.Proxy
}

func (cp *compatibleProvider) Proxies() []C.Proxy {
	return cp.proxies
}

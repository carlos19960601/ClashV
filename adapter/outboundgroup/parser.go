package outboundgroup

import (
	"errors"
	"fmt"

	"github.com/carlos19960601/ClashV/adapter/outbound"
	"github.com/carlos19960601/ClashV/common/structure"
	C "github.com/carlos19960601/ClashV/constant"
)

var (
	errFormat = errors.New("format error")
	errType   = errors.New("unsupported type")
)

type GroupCommonOption struct {
	outbound.BasicOption
	Name    string   `group:"name"`
	Type    string   `group:"type"`
	Proxies []string `group:"proxies,omitempty"`
	Lazy    bool     `group:"lazy,omitempty"`
}

func ParseProxyGroup(config map[string]any, proxyMap map[string]C.Proxy) (C.ProxyAdapter, error) {
	decoder := structure.NewDecoder(structure.Option{TagName: "group", WeaklyTypedInput: true})

	groupOption := &GroupCommonOption{
		Lazy: true,
	}

	if err := decoder.Decode(config, groupOption); err != nil {
		return nil, errFormat
	}

	if groupOption.Type == "" || groupOption.Name == "" {
		return nil, errFormat
	}

	groupName := groupOption.Name

	if len(groupOption.Proxies) != 0 {
		ps, err := getProxies(proxyMap, groupOption.Proxies)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", groupName, err)
		}

		pd, err := provider.NewCompatibelProvider(groupName, ps, hc)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", groupName, err)
		}
		providers = append([]types.ProxyProvider{pd}, providers...)
		providersMap[groupName] = pd
	}

	var group C.ProxyAdapter
	switch groupOption.Type {
	case "select":
		group = NewSelector(groupOption)
	case "fallback":
		group = NewFallback(groupOption)
	default:
		return nil, fmt.Errorf("%w: %s", errType, groupOption.Type)
	}

	return group, nil
}

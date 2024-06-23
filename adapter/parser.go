package adapter

import (
	"fmt"

	"github.com/carlos19960601/ClashV/adapter/outbound"
	"github.com/carlos19960601/ClashV/common/structure"
	C "github.com/carlos19960601/ClashV/constant"
)

func ParseProxy(mapping map[string]any) (C.Proxy, error) {
	decoder := structure.NewDecoder(structure.Option{TagName: "proxy", WeaklyTypedInput: true, KeyReplacer: structure.DefaultKeyReplacer})
	proxyType, existType := mapping["type"].(string)
	if !existType {
		return nil, fmt.Errorf("缺失type")
	}

	var (
		proxy C.ProxyAdapter
		err   error
	)

	switch proxyType {
	case "ss":
	case "ssr":
	case "socks5":
	case "trojan":
		trojanOption := &outbound.TrojanOption{}
		err = decoder.Decode(mapping, trojanOption)
		if err != nil {
			break
		}
		proxy, err = outbound.NewTrojan(*trojanOption)
	default:
		return nil, fmt.Errorf("不支持的proxy类型: %s", proxyType)
	}

	if err != nil {
		return nil, err
	}

	return NewProxy(proxy), nil
}

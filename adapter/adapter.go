package adapter

import (
	"github.com/carlos19960601/ClashV/common/atomic"
	C "github.com/carlos19960601/ClashV/constant"
)

type Proxy struct {
	C.ProxyAdapter
	alive atomic.Bool
}

func NewProxy(adapter C.ProxyAdapter) *Proxy {
	return &Proxy{
		ProxyAdapter: adapter,
		alive:        atomic.NewBool(true),
	}
}

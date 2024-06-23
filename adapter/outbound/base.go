package outbound

import (
	"encoding/json"

	"github.com/carlos19960601/ClashV/common/utils"
	C "github.com/carlos19960601/ClashV/constant"
)

type Base struct {
	name string
	addr string
	tp   C.AdapterType
	udp  bool
	id   string
}

// Id implements C.ProxyAdapter
func (b *Base) Id() string {
	if b.id == "" {
		b.id = utils.NewUUIDV6().String()
	}

	return b.id
}

// Name implements C.ProxyAdapter
func (b *Base) Name() string {
	return b.name
}

// Type implements C.ProxyAdapter
func (b *Base) Type() C.AdapterType {
	return b.tp
}

// Addr implements C.ProxyAdapter
func (b *Base) Addr() string {
	return b.addr
}

// SupportUDP implements C.ProxyAdapter
func (b *Base) SupportUDP() bool {
	return b.udp
}

// MarshalJSON implements C.ProxyAdapter
func (b *Base) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type": b.Type().String(),
		"id":   b.Id(),
	})
}

package context

import (
	"net"

	"github.com/gofrs/uuid/v5"

	N "github.com/carlos19960601/ClashV/common/net"
	"github.com/carlos19960601/ClashV/common/utils"
	C "github.com/carlos19960601/ClashV/constant"
)

type ConnContext struct {
	id       uuid.UUID
	metadata *C.Metadata
	conn     *N.BufferedConn
}

func NewConnContext(conn net.Conn, metadata *C.Metadata) *ConnContext {
	return &ConnContext{
		id:       utils.NewUUIDV4(),
		metadata: metadata,
		conn:     N.NewBufferedConn(conn),
	}
}

// ID implement C.ConnContext ID
func (c *ConnContext) ID() uuid.UUID {
	return c.id
}

// Metadata implement C.ConnContext Metadata
func (c *ConnContext) Metadata() *C.Metadata {
	return c.metadata
}

// Conn implement C.ConnContext Conn
func (c *ConnContext) Conn() *N.BufferedConn {
	return c.conn
}

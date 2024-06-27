package deadline

import (
	"net"
	"time"

	"github.com/carlos19960601/ClashV/common/atomic"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/network"
)

type Conn struct {
	network.ExtendedConn
	deadline atomic.TypedValue[time.Time]
	inRead   atomic.Bool
}

func IsConn(conn any) bool {
	_, ok := conn.(*Conn)
	return ok
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		ExtendedConn: bufio.NewExtendedConn(conn),
	}

	return c
}

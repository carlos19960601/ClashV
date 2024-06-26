package deadline

import (
	"time"

	"github.com/carlos19960601/ClashV/common/atomic"
	"github.com/sagernet/sing/common/network"
)

type Conn struct {
	network.ExtendedConn
	deadline atomic.TypedValue[time.Time]
	inRead   atomic.Bool
}

package tunnel

import (
	"net"

	C "github.com/carlos19960601/ClashV/constant"
)

var (
	status    = newAtomicStatus(Suspend)
	listeners = make(map[string]C.InboundListener)
)

func OnSuspend() {
	status.Store(Suspend)
}

func OnInnerLoading() {
	status.Store(Inner)
}

func OnRunning() {
	status.Store(Running)
}

type tunnel struct{}

var Tunnel C.Tunnel = tunnel{}

// HandleTCPConn implements C.Tunnel
func (t tunnel) HandleTCPConn(conn net.Conn, metadata *C.Metadata) {

}

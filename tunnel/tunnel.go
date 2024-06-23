package tunnel

import (
	"net"

	C "github.com/carlos19960601/ClashV/constant"
)

var (
	status = newAtomicStatus(Suspend)
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

func (t tunnel) HandleTCPConn(conn net.Conn, metadata *C.Metadata) {

}

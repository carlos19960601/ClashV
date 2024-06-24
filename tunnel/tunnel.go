package tunnel

import (
	"net"

	C "github.com/carlos19960601/ClashV/constant"
	icontext "github.com/carlos19960601/ClashV/context"
	"github.com/carlos19960601/ClashV/log"
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
	connCtx := icontext.NewConnContext(conn, metadata)
	handleTCPConn(connCtx)
}

func handleTCPConn(connCtx C.ConnContext) {
	if !isHandle(connCtx.Metadata().Type) {
		_ = connCtx.Conn().Close()
		return
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(connCtx.Conn())

	metadata := connCtx.Metadata()
	if !metadata.Valid() {
		log.Warnln("[Metadata] 无效: %#v", metadata)
		return
	}
}

func isHandle(t C.Type) bool {
	status := status.Load()
	return status == Running || (status == Inner && t == C.INNER)
}

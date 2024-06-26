package loopback

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/puzpuzpuz/xsync/v3"

	C "github.com/carlos19960601/ClashV/constant"
)

var ErrReject = errors.New("reject loopback connection")

type Detector struct {
	connMap       *xsync.MapOf[netip.AddrPort, struct{}]
	packetConnMap *xsync.MapOf[uint16, struct{}]
}

func NewDetector() *Detector {
	return &Detector{
		connMap:       xsync.NewMapOf[netip.AddrPort, struct{}](),
		packetConnMap: xsync.NewMapOf[uint16, struct{}](),
	}
}

func (l *Detector) CheckConn(metadata *C.Metadata) error {
	connAddr := metadata.SourceAddrPort()
	if !connAddr.IsValid() {
		return nil
	}

	if _, ok := l.connMap.Load(connAddr); ok {
		return fmt.Errorf("%w to: %s", ErrReject, metadata.RemoteAddress())
	}

	return nil
}

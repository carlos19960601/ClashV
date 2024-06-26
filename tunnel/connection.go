package tunnel

import (
	"net"

	N "github.com/carlos19960601/ClashV/common/net"
)

func handleSocket(inbound, outbound net.Conn) {
	N.Relay(inbound, outbound)
}

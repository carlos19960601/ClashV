package constant

import "net"

type Tunnel interface {
	HandleTCPConn(conn net.Conn, metadata *Metadata)
	// HandleUDPPacket(packet UDPPacket, metadata *Metadata)
	// NatTable() NatTable
}

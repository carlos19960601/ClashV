package net

import (
	"context"
	"net"

	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/bufio/deadline"
	"github.com/sagernet/sing/common/network"
)

type ExtendedConn = network.ExtendedConn

var NewExtendedConn = bufio.NewExtendedConn

func Relay(leftConn, rightConn net.Conn) {
	_ = bufio.CopyConn(context.TODO(), leftConn, rightConn)
}

func NewDeadlineConn(conn net.Conn) ExtendedConn {
	return deadline.NewConn(conn)
}

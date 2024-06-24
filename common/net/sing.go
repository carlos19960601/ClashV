package net

import (
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/network"
)

type ExtendedConn = network.ExtendedConn

var NewExtendedConn = bufio.NewExtendedConn

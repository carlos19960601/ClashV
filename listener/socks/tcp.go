package socks

import (
	"net"

	"github.com/carlos19960601/ClashV/adapter/inbound"
	N "github.com/carlos19960601/ClashV/common/net"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/transport/socks4"
	"github.com/carlos19960601/ClashV/transport/socks5"
)

type Listener struct {
	listener net.Listener
	addr     string
	closed   bool
}

func New(addr string, tunnel C.Tunnel, additions ...inbound.Addition) (*Listener, error) {
	if len(additions) == 0 {
		additions = []inbound.Addition{}
	}

	l, err := inbound.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	sl := &Listener{
		listener: l,
		addr:     addr,
	}

	go func() {
		for {
			conn, err := sl.listener.Accept()
			if err != nil {
				if sl.closed {
					break
				}
				continue
			}

			go handleSocks(conn, tunnel, additions...)
		}
	}()

	return sl, nil
}

func handleSocks(conn net.Conn, tunnel C.Tunnel, additions ...inbound.Addition) {
	bufConn := N.NewBufferedConn(conn)
	head, err := bufConn.Peek(1)
	if err != nil {
		conn.Close()
		return
	}

	switch head[0] {
	case socks4.Version:
		HandleSocks4(bufConn, tunnel, additions...)
	case socks5.Version:
		HandleSocks5(bufConn, tunnel, additions...)
	default:
		conn.Close()
	}
}

func HandleSocks4(conn net.Conn, tunnel C.Tunnel, additions ...inbound.Addition) {

}

func HandleSocks5(conn net.Conn, tunnel C.Tunnel, additions ...inbound.Addition) {

}

// Close implements C.Listener
func (l *Listener) Close() error {
	l.closed = true
	return l.listener.Close()
}

// Address implements C.Listener
func (l *Listener) Address() string {
	return l.listener.Addr().String()
}

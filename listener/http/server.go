package http

import (
	"net"

	"github.com/carlos19960601/ClashV/adapter/inbound"
	"github.com/carlos19960601/ClashV/common/lru"
	C "github.com/carlos19960601/ClashV/constant"
)

type Listener struct {
	listener net.Listener
	addr     string
	closed   bool
}

func New(addr string, tunnel C.Tunnel, additions ...inbound.Addition) (*Listener, error) {
	return NewWithAuthenticate(addr, tunnel, true, additions...)
}

func NewWithAuthenticate(addr string, tunnel C.Tunnel, authenticate bool, additions ...inbound.Addition) (*Listener, error) {
	if len(additions) == 0 {
		additions = []inbound.Addition{}
	}

	l, err := inbound.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	var c *lru.LruCache[string, bool]
	if authenticate {
		c = lru.New(lru.WithAge[string, bool](30))
	}

	hl := &Listener{
		listener: l,
		addr:     addr,
	}

	go func() {
		for {
			conn, err := hl.listener.Accept()
			if err != nil {
				if hl.closed {
					break
				}
				continue
			}

			go HandleConn(conn, tunnel, c, additions...)
		}
	}()

	return hl, nil
}

func (l *Listener) RawAddress() string {
	return l.addr
}

func (l *Listener) Address() string {
	return l.listener.Addr().String()
}

func (l *Listener) Close() error {
	l.closed = true
	return l.listener.Close()
}

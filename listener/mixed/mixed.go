package mixed

import (
	"net"

	"github.com/carlos19960601/ClashV/adapter/inbound"
	"github.com/carlos19960601/ClashV/common/lru"
	N "github.com/carlos19960601/ClashV/common/net"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/listener/http"
	"github.com/carlos19960601/ClashV/listener/socks"
	"github.com/carlos19960601/ClashV/transport/socks4"
	"github.com/carlos19960601/ClashV/transport/socks5"
)

type Listener struct {
	listener net.Listener
	addr     string
	closed   bool
	cache    *lru.LruCache[string, bool]
}

func New(addr string, tunnel C.Tunnel, additions ...inbound.Addition) (*Listener, error) {
	if len(additions) == 0 {
		additions = []inbound.Addition{}
	}

	l, err := inbound.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ml := &Listener{
		listener: l,
		addr:     addr,
		cache:    lru.New(lru.WithAge[string, bool](30)),
	}

	go func() {
		for {
			conn, err := ml.listener.Accept()
			if err != nil {
				if ml.closed {
					break
				}
				continue
			}

			go handleConn(conn, tunnel, ml.cache, additions...)
		}
	}()

	return ml, nil
}

// Address implements C.Listener
func (l *Listener) Address() string {
	return l.listener.Addr().String()
}

// Close implements C.Listener
func (l *Listener) Close() error {
	l.closed = true
	return l.listener.Close()
}

func handleConn(conn net.Conn, tunnel C.Tunnel, cache *lru.LruCache[string, bool], additions ...inbound.Addition) {
	bufConn := N.NewBufferedConn(conn)
	head, err := bufConn.Peek(1)
	if err != nil {
		return
	}

	switch head[0] {
	case socks4.Version:
		socks.HandleSocks4(bufConn, tunnel, additions...)
	case socks5.Version:
		socks.HandleSocks5(bufConn, tunnel, additions...)
	default:
		http.HandleConn(bufConn, tunnel, cache, additions...)
	}
}

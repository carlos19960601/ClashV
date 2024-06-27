package callback

import (
	"sync"

	C "github.com/carlos19960601/ClashV/constant"
)

type closeCallbackConn struct {
	C.Conn
	coloseFunc func()
	coloseOnce sync.Once
}

func NewCloseCallbackConn(conn C.Conn, callback func()) C.Conn {
	return &closeCallbackConn{Conn: conn, coloseFunc: callback}
}

func (w *closeCallbackConn) Colse() error {
	w.coloseOnce.Do(w.coloseFunc)
	return w.Conn.Close()
}

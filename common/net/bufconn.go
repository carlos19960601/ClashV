package net

import (
	"bufio"
	"net"
)

var _ ExtendedConn = (*BufferedConn)(nil)

type BufferedConn struct {
	r *bufio.Reader
	ExtendedConn
	peeked bool
}

func NewBufferedConn(c net.Conn) *BufferedConn {
	if bc, ok := c.(*BufferedConn); ok {
		return bc
	}

	return &BufferedConn{bufio.NewReader(c), NewExtendedConn(c), false}
}

func (c *BufferedConn) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

func (c *BufferedConn) Reader() *bufio.Reader {
	return c.r
}

// Peek returns the next n bytes without advancing the reader.
func (c *BufferedConn) Peek(n int) ([]byte, error) {
	c.peeked = true
	return c.r.Peek(n)
}

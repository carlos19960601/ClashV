package common

import "errors"

var (
	errPayload = errors.New("payloadRule error")
	noResolve  = "no-resolve"
)

type Base struct {
}

func (n *Base) ShouldResolveIP() bool {
	return false
}

func (b *Base) ShouldFindProcess() bool {
	return false
}

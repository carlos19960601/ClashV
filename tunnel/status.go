package tunnel

import "sync/atomic"

type TunnelStatus int

var StatusMapping = map[string]TunnelStatus{
	Suspend.String(): Suspend,
	Inner.String():   Inner,
	Running.String(): Running,
}

const (
	Suspend TunnelStatus = iota
	Inner
	Running
)

func (s TunnelStatus) String() string {
	switch s {
	case Suspend:
		return "suspend"
	case Inner:
		return "inner"
	case Running:
		return "running"
	default:
		return "Unknown"
	}
}

type AtomicStatus struct {
	value atomic.Int32
}

func (a *AtomicStatus) Store(s TunnelStatus) {
	a.value.Store(int32(s))
}

func (a *AtomicStatus) Load() TunnelStatus {
	return TunnelStatus(a.value.Load())
}

func (a *AtomicStatus) String() string {
	return a.Load().String()
}

func newAtomicStatus(s TunnelStatus) *AtomicStatus {
	a := &AtomicStatus{}
	a.Store(s)
	return a
}

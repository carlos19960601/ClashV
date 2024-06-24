package atomic

import "sync/atomic"

type Bool struct {
	atomic.Bool
}

func NewBool(val bool) (i Bool) {
	i.Store(val)
	return
}



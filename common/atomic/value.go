package atomic

import "sync/atomic"

type TypeValue[T any] struct {
	_     noCopy
	value atomic.Value
}

func DefaultValue[T any]() T {
	var defaultValue T
	return defaultValue
}

// tValue is a struct with determined type to resolve atomic.Value usages with interface types
// https://github.com/golang/go/issues/22550
type tValue[T any] struct {
	value T
}

func (t *TypeValue[T]) Load() T {
	value := t.value.Load()
	if value == nil {
		return DefaultValue[T]()
	}

	return value.(tValue[T]).value
}

func (t *TypeValue[T]) Store(value T) {
	t.value.Store(tValue[T]{value})
}

type noCopy struct{}

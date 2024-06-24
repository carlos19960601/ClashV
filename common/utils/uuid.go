package utils

import (
	"github.com/gofrs/uuid/v5"
	"github.com/zhangyunhao116/fastrand"
)

type fastRandReader struct{}

func (r fastRandReader) Read(p []byte) (int, error) {
	return fastrand.Read(p)
}

var UnsafeUUIDGenerator = uuid.NewGenWithOptions(uuid.WithRandomReader(fastRandReader{}))

func NewUUIDV6() uuid.UUID {
	u, _ := UnsafeUUIDGenerator.NewV6() // fastrand.Read wouldn't cause error, so ignore err is safe
	return u
}

func NewUUIDV4() uuid.UUID {
	u, _ := UnsafeUUIDGenerator.NewV4() // fastrand.Read wouldn't cause error, so ignore err is safe
	return u
}

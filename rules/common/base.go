package common

import "errors"

var (
	errPayload = errors.New("payloadRule error")
	noResolve  = "no-resolve"
)

type Base struct {
}

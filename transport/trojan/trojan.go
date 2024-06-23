package trojan

import (
	"crypto/sha256"
	"encoding/hex"
)

type Option struct {
	Password       string
	ServerName     string
	SkipCertVerify bool
}

type Trojan struct {
	option      *Option
	hexPassword []byte
}

func New(option *Option) *Trojan {
	return &Trojan{option, hexSha224([]byte(option.Password))}
}

func hexSha224(data []byte) []byte {
	buf := make([]byte, 56)
	hash := sha256.Sum224(data)
	hex.Encode(buf, hash[:])
	return buf
}

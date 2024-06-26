package common

import (
	"strings"

	"golang.org/x/net/idna"

	C "github.com/carlos19960601/ClashV/constant"
)

type DomainSuffix struct {
	*Base
	suffix  string
	adapter string
}

func (ds *DomainSuffix) RuleType() C.RuleType {
	return C.DomainSuffix
}

func (ds *DomainSuffix) Adapter() string {
	return ds.adapter
}

func (ds *DomainSuffix) Payload() string {
	return ds.suffix
}

func (ds *DomainSuffix) Match(metadata *C.Metadata) (bool, string) {
	domain := metadata.RuleHost()
	return strings.HasSuffix(domain, "."+ds.suffix) || domain == ds.suffix, ds.adapter
}

func NewDomainSuffix(suffix string, adapter string) *DomainSuffix {
	punycode, _ := idna.ToASCII(strings.ToLower(suffix))
	return &DomainSuffix{
		Base:    &Base{},
		suffix:  punycode,
		adapter: adapter,
	}
}

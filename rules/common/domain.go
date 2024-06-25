package common

import (
	"strings"

	C "github.com/carlos19960601/ClashV/constant"
	"golang.org/x/net/idna"
)

type Domain struct {
	*Base
	domain  string
	adapter string
}

func (d *Domain) RuleType() C.RuleType {
	return C.Domain
}

func (d *Domain) Adapter() string {
	return d.adapter
}

func (d *Domain) Payload() string {
	return d.domain
}

func (d *Domain) Match(metadata *C.Metadata) (bool, string) {
	return metadata.RuleHost() == d.domain, d.adapter
}

func NewDomain(domain string, adapter string) *Domain {
	punycode, _ := idna.ToASCII(strings.ToLower(domain))
	return &Domain{
		Base:    &Base{},
		domain:  punycode,
		adapter: adapter,
	}
}

package inbound

import (
	"net"

	C "github.com/carlos19960601/ClashV/constant"
)

type Addition func(metadata *C.Metadata)

func ApplyAdditions(metadata *C.Metadata, additions ...Addition) {
	for _, addition := range additions {
		addition(metadata)
	}
}

func WithSrcAddr(addr net.Addr) Addition {
	return func(metadata *C.Metadata) {
		m := C.Metadata{}
		if err := m.SetRemoteAddr(addr); err == nil {
			metadata.SrcIP = m.DstIP
			metadata.SrcPort = m.DstPort
		}
	}
}

func WithInAddr(addr net.Addr) Addition {
	return func(metadata *C.Metadata) {
		m := &C.Metadata{}
		if err := m.SetRemoteAddr(addr); err == nil {
			metadata.InIP = m.DstIP
			metadata.InPort = m.DstPort
		}
	}
}

func WithInName(name string) Addition {
	return func(metadata *C.Metadata) {
		metadata.InName = name
	}
}

func WithSpecialRules(specialRules string) Addition {
	return func(metadata *C.Metadata) {
		metadata.SpecialRules = specialRules
	}
}

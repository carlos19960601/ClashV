package provider

const (
	Proxy ProviderType = iota
	Rule
)

type ProviderType int

func (pt ProviderType) String() string {
	switch pt {
	case Proxy:
		return "Proxy"
	case Rule:
		return "Rule"
	default:
		return "Unknown"
	}
}

type Provider interface {
	Name() string
}

type ProxyProvider interface {
	Provider
}



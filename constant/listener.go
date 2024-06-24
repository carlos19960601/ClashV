package constant

type InboundListener interface {
	Name() string
	Listen(tunnel Tunnel) error
	Close() error
	Address() string
	RawAddress() string
	Config() InboundConfig
}

type InboundConfig interface {
	Name() string
	Equal(config InboundConfig) bool
}

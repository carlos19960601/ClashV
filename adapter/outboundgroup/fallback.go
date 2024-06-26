package outboundgroup

type Fallback struct {
	*GroupBase
}

func NewFallback(option *GroupCommonOption) *Fallback {
	return &Fallback{}
}

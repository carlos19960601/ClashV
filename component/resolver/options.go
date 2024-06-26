package resolver

import "github.com/carlos19960601/ClashV/component/resolver"

type option struct {
	resolver resolver.Resolver
}

type Option func(opt *option)

func WithResolver(r resolver.Resolver) Option {
	return func(opt *option) {
		opt.resolver = r
	}
}

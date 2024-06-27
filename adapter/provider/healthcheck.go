package provider

import (
	"sync"
	"time"

	C "github.com/carlos19960601/ClashV/constant"
)

type HealthCheckOption struct {
	URL      string
	Interval uint
}

type HealthCheck struct {
	url      string
	mu       sync.Mutex
	proxies  []C.Proxy
	interval time.Duration
	timeout  time.Duration
}

func NewHealthCheck(proxies []C.Proxy, url string, timeout uint, interval uint, lazy bool) *HealthCheck {
	if url == "" {
		interval = 0
	}

	if timeout == 0 {
		timeout = 5000
	}

	return &HealthCheck{
		proxies:  proxies,
		url:      url,
		timeout:  time.Duration(timeout) * time.Millisecond,
		interval: time.Duration(interval) * time.Second,
	}
}

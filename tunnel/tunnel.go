package tunnel

import (
	"context"
	"net"
	"sync"

	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/constant/provider"
	icontext "github.com/carlos19960601/ClashV/context"
	"github.com/carlos19960601/ClashV/log"
)

var (
	status    = newAtomicStatus(Suspend)
	listeners = make(map[string]C.InboundListener)
	proxies   = make(map[string]C.Proxy)
	providers map[string]provider.ProxyProvider
	configMux sync.RWMutex

	// 出站规则
	mode = Rule
)

func OnSuspend() {
	status.Store(Suspend)
}

func OnInnerLoading() {
	status.Store(Inner)
}

func OnRunning() {
	status.Store(Running)
}

func Mode() TunnelMode {
	return mode
}

func SetMode(m TunnelMode) {
	mode = m
}

func UpdateProxies(newProxies map[string]C.Proxy, newProviders map[string]provider.ProxyProvider) {
	configMux.Lock()
	defer configMux.Unlock()

	proxies = newProxies
	providers = newProviders
}

func updateListeners(newListeners map[string]C.InboundListener) {
	configMux.Lock()
	defer configMux.Unlock()

	listeners = newListeners
}

type tunnel struct{}

var Tunnel C.Tunnel = tunnel{}

// HandleTCPConn implements C.Tunnel
func (t tunnel) HandleTCPConn(conn net.Conn, metadata *C.Metadata) {
	connCtx := icontext.NewConnContext(conn, metadata)
	handleTCPConn(connCtx)
}

func handleTCPConn(connCtx C.ConnContext) {
	if !isHandle(connCtx.Metadata().Type) {
		_ = connCtx.Conn().Close()
		return
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(connCtx.Conn())

	metadata := connCtx.Metadata()
	if !metadata.Valid() {
		log.Warnln("[Metadata] 无效: %#v", metadata)
		return
	}

	proxy, rule, err := resolveMetadata(metadata)
	if err != nil {
		log.Warnln("[Metadata] 解析失败: %s", err.Error())
		return
	}

	dialMetadata := metadata
	if len(metadata.Host) > 0 {

	}

	ctx, cancel := context.WithTimeout(context.Background(), C.DefaultTCPTimeout)
	defer cancel()
	remoteConn, err := retry(ctx, func(ctx context.Context) (remoteConn C.Conn, err error) {
		remoteConn, err = proxy.DialContext(ctx, dialMetadata)

		return
	}, func(err error) {
		if rule == nil {

		}
	})

	defer func(remoteConn C.Conn) {
		_ = remoteConn.Close()
	}(remoteConn)

}

func resolveMetadata(metadata *C.Metadata) (proxy C.Proxy, rule C.Rule, err error) {
	switch mode {
	case Direct:
		proxy = proxies["DIRECT"]
	case Global:
		proxy = proxies["GLOBAL"]
	// Rule
	default:
		proxy, rule, err = match(metadata)
	}

	return
}

func match(metadata *C.Metadata) (C.Proxy, C.Rule, error) {
	configMux.RLock()
	defer configMux.RUnlock()

	for range getRules(metadata) {

	}

	return nil, nil, nil
}

func isHandle(t C.Type) bool {
	status := status.Load()
	return status == Running || (status == Inner && t == C.INNER)
}

func getRules(metadata *C.Metadata) []C.Rule {
	return nil
}

func retry[T any](ctx context.Context, ft func(context.Context) (T, error), fe func(err error)) (t T, err error) {
	for i := 0; i < 10; i++ {
		t, err = ft(ctx)
		if err != nil {
			if fe != nil {
				fe(err)
			}

		} else {
			break
		}
	}

	return
}

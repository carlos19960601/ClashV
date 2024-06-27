package tunnel

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/carlos19960601/ClashV/component/resolver"
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
	rules     []C.Rule
	subRules  map[string][]C.Rule

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

func UpdateRules(newRules []C.Rule, newSubRules map[string][]C.Rule) {
	configMux.Lock()
	defer configMux.Unlock()
	rules = newRules
	subRules = newSubRules
}

func UpdateListeners(newListeners map[string]C.InboundListener) {
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

	conn := connCtx.Conn()
	proxy, rule, err := resolveMetadata(metadata)
	if err != nil {
		log.Warnln("[Metadata] 解析失败: %s", err.Error())
		return
	}

	dialMetadata := metadata
	if len(metadata.Host) > 0 {
		if node, ok := resolver.DefaultHosts.Search(metadata.Host, false); ok {
			dstIp, _ := node.RandIP()
			dialMetadata.DstIP = dstIp
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), C.DefaultTCPTimeout)
	defer cancel()
	remoteConn, err := retry(ctx, func(ctx context.Context) (remoteConn C.Conn, err error) {
		remoteConn, err = proxy.DialContext(ctx, dialMetadata)
		if err != nil {
			return nil, err
		}

		return
	}, func(err error) {
		if rule == nil {

		}
	})

	defer func(remoteConn C.Conn) {
		if remoteConn != nil {
			_ = remoteConn.Close()
		}
	}(remoteConn)

	switch true {
	case rule != nil:
		if rule.Payload() != "" {
			log.Infoln("[TCP] %s --> %s 匹配 %s using %s", metadata.SourceDetail(), metadata.RemoteAddress(), fmt.Sprintf("%s(%s)", rule.RuleType().String(), rule.Payload()), remoteConn.Chains().String())
		} else {
			log.Infoln("[TCP] %s --> %s 匹配 %s using %s", metadata.SourceDetail(), metadata.RemoteAddress(), rule.RuleType().String(), remoteConn.Chains().String())
		}
	case mode == Global:
		log.Infoln("[TCP] %s --> % 使用 GLOBAL", metadata.SourceDetail(), metadata.RemoteAddress())
	default:
		log.Infoln(
			"[TCP] %s --> %s 没有匹配到规则, 使用DIRECT",
			metadata.SourceDetail(),
			metadata.RemoteAddress(),
		)
	}

	handleSocket(conn, remoteConn)
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

	var (
		resolved bool
	)

	if node, ok := resolver.DefaultHosts.Search(metadata.Host, false); ok {
		metadata.DstIP, _ = node.RandIP()
		resolved = true
	}

	for _, rule := range getRules(metadata) {
		if !resolved && shouldResolveIP(rule, metadata) {
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), resolver.DefaultDNSTimeout)
				defer cancel()

				ip, err := resolver.ResolveIP(ctx, metadata.Host)
				if err != nil {
					log.Debugln("[DNS] resolve %s error: %s", metadata.Host, err.Error())
				} else {
					log.Debugln("[DNS] %s --> %s", metadata.Host, ip.String())
					metadata.DstIP = ip
				}

				resolved = true
			}()
		}

		// ada是 "一元机场" 等
		if matched, ada := rule.Match(metadata); matched {
			adapter, ok := proxies[ada]
			if !ok {
				continue
			}

			passed := false
			for adapter := adapter; adapter != nil; adapter = adapter.Unwrap(metadata, false) {
				if adapter.Type() == C.Pass {
					passed = true
					break
				}
			}

			if passed {
				log.Debugln("%s match Pass rule", adapter.Name())
				continue
			}

			return adapter, rule, nil
		}
	}

	return proxies["DIRECT"], nil, nil
}

func isHandle(t C.Type) bool {
	status := status.Load()
	return status == Running || (status == Inner && t == C.INNER)
}

func shouldResolveIP(rule C.Rule, metadata *C.Metadata) bool {
	return rule.ShouldResolveIP() && metadata.Host != "" && !metadata.DstIP.IsValid()
}

func getRules(metadata *C.Metadata) []C.Rule {
	if sr, ok := subRules[metadata.SpecialRules]; ok {
		log.Debugln("[Rule] 使用 %s 规则", metadata.SpecialRules)
		return sr
	} else {
		log.Debugln("[Rule] 使用默认规则")
		return rules
	}
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

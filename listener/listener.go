package listener

import (
	"fmt"
	"net"
	"sync"

	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/listener/http"
	"github.com/carlos19960601/ClashV/listener/mixed"
	"github.com/carlos19960601/ClashV/listener/socks"
	"github.com/carlos19960601/ClashV/log"
)

var (
	allowLan    = false
	bindAddress = "*"

	socksListener *socks.Listener
	httpListener  *http.Listener
	mixedListener *mixed.Listener

	inboundListeners = map[string]C.InboundListener{}

	socksMux   sync.Mutex
	inboundMux sync.Mutex
	httpMux    sync.Mutex
	mixedMux   sync.Mutex
)

func SetAllowLan(al bool) {
	allowLan = al
}

func SetBindAddress(host string) {
	bindAddress = host
}

func PatchInboundListeners(newListenerMap map[string]C.InboundListener, tunnel C.Tunnel, dropOld bool) {
	inboundMux.Lock()
	defer inboundMux.Unlock()

	for name, newListener := range newListenerMap {
		if oldListener, ok := inboundListeners[name]; ok {
			if !oldListener.Config().Equal(newListener.Config()) {
				_ = oldListener.Close()
			} else {
				continue
			}
		}

		inboundListeners[name] = newListener
	}
}

func ReCreateHTTP(port int, tunnel C.Tunnel) {
	httpMux.Lock()
	defer httpMux.Unlock()

	var err error
	defer func() {
		if err != nil {
			log.Errorln("启动HTTP服务失败: %s", err.Error())
		}
	}()

	addr := genAddr(bindAddress, port, allowLan)

	if httpListener != nil {
		if httpListener.RawAddress() == addr {
			return
		}

		httpListener.Close()
		httpListener = nil
	}

	if portIsZero(addr) {
		return
	}

	httpListener, err = http.New(addr, tunnel)
	if err != nil {
		log.Errorln("启动HTTP服务失败: %s", err.Error())
		return
	}

	log.Infoln("HTTP代理监听在: %s", httpListener.Address())
}

func ReCreateSocks(port int, tunnel C.Tunnel) {
	socksMux.Lock()
	defer socksMux.Unlock()

	var err error
	defer func() {
		if err != nil {
			log.Errorln("启动Socks服务失败: %s", err.Error())
		}
	}()

	addr := genAddr(bindAddress, port, allowLan)

	if portIsZero(addr) {
		return
	}

	tcpListener, err := socks.New(addr, tunnel)
	if err != nil {
		tcpListener.Close()
		return
	}

	socksListener = tcpListener

	log.Infoln("SOCKS 代理监听在: %s", socksListener.Address())
}

func ReCreateMixed(port int, tunnel C.Tunnel) {
	mixedMux.Lock()
	defer mixedMux.Unlock()

	var err error
	defer func() {
		if err != nil {
			log.Errorln("Start Mixed(http+socks) server error: %s", err.Error())
		}
	}()

	addr := genAddr(bindAddress, port, allowLan)

	if portIsZero(addr) {
		return
	}
	mixedListener, err = mixed.New(addr, tunnel)
	if err != nil {
		return
	}

	log.Infoln("Mixed(http+socks) 代理监听在: %s", mixedListener.Address())
}

func genAddr(host string, port int, allowLan bool) string {
	if allowLan {
		if host == "*" {
			return fmt.Sprintf(":%d", port)
		}
		return fmt.Sprintf("%s:%d", host, port)
	}

	return fmt.Sprintf("127.0.0.1:%d", port)
}

func portIsZero(addr string) bool {
	_, port, err := net.SplitHostPort(addr)
	if port == "0" || port == "" || err != nil {
		return true
	}

	return false
}

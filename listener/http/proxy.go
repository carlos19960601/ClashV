package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/carlos19960601/ClashV/adapter/inbound"
	"github.com/carlos19960601/ClashV/common/lru"
	N "github.com/carlos19960601/ClashV/common/net"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/log"
)

func HandleConn(c net.Conn, tunnel C.Tunnel, cache *lru.LruCache[string, bool], additions ...inbound.Addition) {
	client := newClient(c, tunnel, additions...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	peekMutex := sync.Mutex{}
	conn := N.NewBufferedConn(c)

	keepAlive := true
	trusted := cache == nil

	if !trusted {

	}

	if trusted {

	}

	for keepAlive {
		peekMutex.Lock()
		request, err := ReadRequest(conn.Reader())
		peekMutex.Unlock()
		if err != nil {
			break
		}

		request.RemoteAddr = conn.RemoteAddr().String()
		keepAlive = strings.TrimSpace(strings.ToLower(request.Header.Get("Proxy-Connection"))) == "keep-alive"

		var resp *http.Response

		// 走代理的时候，发送的是Connect的请求
		log.Infoln("request.Method: %s", request.Method)
		if request.Method == http.MethodConnect {
			// 手动返回
			if _, err = fmt.Fprintf(conn, "HTTP/%d.%d %03d %s\r\n\r\n", request.ProtoMajor, request.ProtoMinor, http.StatusOK, "Connection established"); err != nil {
				break
			}

			tunnel.HandleTCPConn(inbound.NewHTTPS(request, conn, additions...))

			return
		}

		host := request.Header.Get("Host")
		if host != "" {
			request.Host = host
		}

		// RequestURI 在http client中不能设置RequestURI，需要去掉
		request.RequestURI = ""

		log.Infoln("Schema: %s, Host: %s", request.URL.Scheme, request.URL.Host)
		if request.URL.Scheme == "" || request.URL.Host == "" {
			resp = responseWith(request, http.StatusBadRequest)
		} else {
			request = request.WithContext(ctx)
			resp, err = client.Do(request)
			if err != nil {
				log.Errorln("请求失败: %s", err.Error())
				resp = responseWith(request, http.StatusBadGateway)
			}
		}

		log.Infoln("%+v", resp)
		err = resp.Write(conn)
		if err != nil {
			break // close connection
		}
	}

	_ = conn.Close()
}

func responseWith(request *http.Request, statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Proto:      request.Proto,
		ProtoMajor: request.ProtoMajor,
		ProtoMinor: request.ProtoMinor,
		Header:     http.Header{},
	}
}

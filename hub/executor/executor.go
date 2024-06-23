package executor

import (
	"fmt"
	"os"
	"sync"

	"github.com/carlos19960601/ClashV/config"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/log"
	"github.com/carlos19960601/ClashV/tunnel"
)

var mux sync.Mutex

func Parse() (*config.Config, error) {
	return ParseWithPath(C.Path.Config())
}

func ParseWithPath(path string) (*config.Config, error) {
	buf, err := readConfig(path)
	if err != nil {
		return nil, err
	}

	return ParseWithBytes(buf)
}

func ParseWithBytes(buf []byte) (*config.Config, error) {
	return config.Parse(buf)
}

func ApplyConfig(cfg *config.Config, force bool) {
	mux.Lock()
	defer mux.Unlock()

	tunnel.OnSuspend()

	tunnel.OnInnerLoading()

	initInnerTcp()
	loadProxyProvider(cfg.Providers)
	tunnel.OnRunning()

	log.SetLevel(cfg.General.LogLevel)
}

func initInnerTcp() {
	inner.New(tunnel.Tunnel)
}

func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("配置文件 %s 是空的", path)
	}

	return data, nil
}

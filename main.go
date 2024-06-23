package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/carlos19960601/ClashV/config"
	C "github.com/carlos19960601/ClashV/constant"
	"github.com/carlos19960601/ClashV/hub"
	"github.com/carlos19960601/ClashV/log"
)

var (
	version            bool
	homeDir            string
	configFile         string
	externalController string
)

func init() {
	flag.StringVar(&homeDir, "d", os.Getenv("CLASH_HOME_DIR"), "设置配置文件路径")
	flag.StringVar(&configFile, "f", os.Getenv("CLASH_CONFIG_FILE"), "指定配置文件")
	flag.StringVar(&externalController, "ext-ctl", os.Getenv("CLASH_OVERRIDE_EXTERNAL_CONTROLLER"), "覆盖external controller地址")
	flag.BoolVar(&version, "v", false, "显示版本信息")
	flag.Parse()
}

func main() {
	if version {
		fmt.Printf("ClashV Meta %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime)
	}

	if homeDir != "" {
		if !filepath.IsAbs(homeDir) {
			currentDir, _ := os.Getwd()
			homeDir = filepath.Join(currentDir, homeDir)
		}
		C.SetHomeDir(homeDir)
	}

	if configFile != "" {
		if !filepath.IsAbs(configFile) {
			currentDir, _ := os.Getwd()
			configFile = filepath.Join(currentDir, configFile)
		}
	} else {
		configFile = filepath.Join(C.Path.HomeDir(), C.Path.Config())
	}
	C.SetConfig(configFile)

	if err := config.Init(C.Path.HomeDir()); err != nil {
		log.Fatalln("初始化配置路径失败: %s", err.Error())
	}

	var options []hub.Option
	if externalController != "" {
		options = append(options, hub.WithExternalController(externalController))
	}

	if err := hub.Parse(options...); err != nil {
		log.Fatalln("解析配置失败: %s", err.Error())
	}

	termSign := make(chan os.Signal, 1)
	hupSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(hupSign, syscall.SIGHUP)

	for {
		select {
		case <-termSign:
			return
		// 重新初始化
		case <-hupSign:

		}
	}

}

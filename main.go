package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	C "github.com/carlos19960601/ClashV/constant"
)

var (
	version    bool
	homeDir    string
	configFile string
)

func init() {
	flag.StringVar(&homeDir, "d", os.Getenv("CLASH_HOME_DIR"), "设置配置文件路径")
	flag.StringVar(&configFile, "f", os.Getenv("CLASH_CONFIG_FILE"), "指定配置文件")
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
}

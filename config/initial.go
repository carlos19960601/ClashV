package config

import (
	"fmt"
	"os"

	C "github.com/carlos19960601/ClashV/constant"
)

func Init(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o777); err != nil {
			return fmt.Errorf("不能创建配置文件目录：%w", err)
		}
	}

	// 初始化config.yaml
	if _, err := os.Stat(C.Path.Config()); os.IsNotExist(err) {
		f, err := os.OpenFile(C.Path.Config(), os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("不能创建配置文件：%s: %s", C.Path.Config(), err.Error())
		}
		f.Write([]byte("mixed-port: 7890"))
		f.Close()
	}
	return nil
}

package config

import (
	"fmt"
	"time"

	"github.com/carlos19960601/ClashV/adapter"
	C "github.com/carlos19960601/ClashV/constant"
	providerTypes "github.com/carlos19960601/ClashV/constant/provider"
	"github.com/carlos19960601/ClashV/log"
	T "github.com/carlos19960601/ClashV/tunnel"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General   *General
	Proxies   map[string]C.Proxy
	Providers map[string]providerTypes.ProxyProvider
}

type RawConfig struct {
	Port        int    `yaml:"port" json:"port"`
	BindAddress string `yaml:"bind-address" json:"bind-address"`
	SocksPort   int    `yaml:"socks-port" json:"socks-port"`
	MixedPort   int    `yaml:"mixed-port" json:"mixed-port"`
	AllowLan    bool   `yaml:"allow-lan" json:"allow-lan"`

	Proxy      []map[string]any `yaml:"proxies"`
	ProxyGroup []map[string]any `yaml:"proxy-groups"`
	Rule       []string         `yaml:"rules"`
}

type General struct {
	Inbound
	Controller

	Mode     T.TunnelMode `json:"mode"`
	LogLevel log.LogLevel `json:"log-level"`
}

type Inbound struct {
	Port      int  `json:"port"`
	SocksPort int  `json:"socks-port"`
	MixedPort int  `json:"mixed-port"`
	AllowLan  bool `json:"allow-lan"`
}

type Controller struct {
	ExternalController string `json:"-"`
	ExternalUI         string `json:"-"`
}

func Parse(buf []byte) (*Config, error) {
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		return nil, err
	}

	return ParseRawConfig(rawCfg)
}

func ParseRawConfig(rawCfg *RawConfig) (*Config, error) {
	config := &Config{}
	startTime := time.Now()

	proxies, providers, err := parseProxies(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Proxies = proxies
	config.Providers = providers

	elapsedTime := time.Since(startTime)
	log.Infoln("初始化配置完成， 耗时: %dms", elapsedTime.Seconds())

	return config, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, providerMap map[string]providerTypes.ProxyProvider, err error) {
	proxies = make(map[string]C.Proxy)
	providerMap = make(map[string]providerTypes.ProxyProvider)
	proxiesConfig := cfg.Proxy
	groupConfig := cfg.ProxyGroup
	providerConfig := cfg.ProxyProvider

	for idx, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			return nil, nil, fmt.Errorf("proxy: %d: %w", idx, err)
		}
	}

	for idx, mapping := range groupConfig {
		groupName, existName := mapping["name"].(string)
		if !existName {
			return nil, nil, fmt.Errorf("代理组 %d: 名称缺失", idx)
		}

	}

}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{
		AllowLan:    false,
		BindAddress: "*",
	}

	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}

	return rawCfg, nil
}

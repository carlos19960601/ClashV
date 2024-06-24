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
	Listeners map[string]C.InboundListener
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

	Listeners []map[string]any `yaml:"listeners"`
}

type General struct {
	Inbound
	Controller

	Mode     T.TunnelMode `json:"mode"`
	LogLevel log.LogLevel `json:"log-level"`
}

type Inbound struct {
	Port        int    `json:"port"`
	SocksPort   int    `json:"socks-port"`
	MixedPort   int    `json:"mixed-port"`
	AllowLan    bool   `json:"allow-lan"`
	BindAddress string `json:"bind-address"`
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

	general, err := parseGeneral(rawCfg)
	if err != nil {
		return nil, err
	}
	config.General = general

	proxies, providers, err := parseProxies(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Proxies = proxies
	config.Providers = providers

	elapsedTime := time.Since(startTime) / time.Millisecond
	log.Infoln("初始化配置完成， 耗时: %dms", elapsedTime)

	return config, nil
}

func parseGeneral(cfg *RawConfig) (*General, error) {
	return &General{
		Inbound: Inbound{
			Port:      cfg.Port,
			MixedPort: cfg.MixedPort,
		},
	}, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, providerMap map[string]providerTypes.ProxyProvider, err error) {
	proxies = make(map[string]C.Proxy)
	providerMap = make(map[string]providerTypes.ProxyProvider)
	proxiesConfig := cfg.Proxy
	groupConfig := cfg.ProxyGroup

	var proxyList []string

	for idx, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			return nil, nil, fmt.Errorf("proxy: %d: %w", idx, err)
		}

		if _, exist := proxies[proxy.Name()]; exist {
			return nil, nil, fmt.Errorf("代理: %s 名称重复", proxy.Name())
		}

		proxies[proxy.Name()] = proxy
		proxyList = append(proxyList, proxy.Name())
	}

	for idx, mapping := range groupConfig {
		groupName, existName := mapping["name"].(string)
		if !existName {
			return nil, nil, fmt.Errorf("代理组 %d: 名称缺失", idx)
		}
		proxyList = append(proxyList, groupName)
	}

	var ps []C.Proxy
	for _, v := range proxyList {
		if proxies[v].Type() == C.Pass {
			continue
		}

		ps = append(ps, proxies[v])
	}

	return proxies, providerMap, nil

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

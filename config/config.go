package config

import (
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/carlos19960601/ClashV/adapter"
	"github.com/carlos19960601/ClashV/adapter/outbound"
	"github.com/carlos19960601/ClashV/adapter/outboundgroup"
	"github.com/carlos19960601/ClashV/adapter/provider"
	"github.com/carlos19960601/ClashV/component/resolver"
	"github.com/carlos19960601/ClashV/component/trie"
	C "github.com/carlos19960601/ClashV/constant"
	providerTypes "github.com/carlos19960601/ClashV/constant/provider"
	"github.com/carlos19960601/ClashV/log"
	R "github.com/carlos19960601/ClashV/rules"
	T "github.com/carlos19960601/ClashV/tunnel"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General   *General
	DNS       *DNS
	Proxies   map[string]C.Proxy
	Providers map[string]providerTypes.ProxyProvider
	Listeners map[string]C.InboundListener
	Hosts     *trie.DomainTrie[resolver.HostValue]
	Rules     []C.Rule
	SubRules  map[string][]C.Rule
}

type RawConfig struct {
	Port        int    `yaml:"port" json:"port"`
	BindAddress string `yaml:"bind-address" json:"bind-address"`
	SocksPort   int    `yaml:"socks-port" json:"socks-port"`
	MixedPort   int    `yaml:"mixed-port" json:"mixed-port"`
	AllowLan    bool   `yaml:"allow-lan" json:"allow-lan"`

	Hosts         map[string]any            `yaml:"hosts" json:"hosts"`
	Proxy         []map[string]any          `yaml:"proxies"`
	ProxyGroup    []map[string]any          `yaml:"proxy-groups"`
	Rule          []string                  `yaml:"rules"`
	SubRules      map[string][]string       `yaml:"sub-rules"`
	DNS           RawDNS                    `yaml:"dns" json:"dns"`
	ProxyProvider map[string]map[string]any `yaml:"proxy-providers"`

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

type DNS struct {
	Enable bool `yaml:"enable"`
}

type RawDNS struct {
	Enable            bool      `yaml:"enable" json:"enable"`
	IPv6              bool      `yaml:"ipv6" json:"ipv6"`
	DefaultNameserver []string  `yaml:"default-nameserver" json:"default-nameserver"`
	EnhancedMode      C.DNSMode `yaml:"enhanced-mode" json:"enhanced-mode"`
}

func ParseRawConfig(rawCfg *RawConfig) (*Config, error) {
	config := &Config{}
	log.Infoln("开始初始化配置")
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

	subRules, err := parseSubRules(rawCfg, proxies)
	if err != nil {
		return nil, err
	}
	config.SubRules = subRules

	rules, err := parseRules(rawCfg.Rule, proxies, subRules, "rules")
	if err != nil {
		return nil, err
	}
	config.Rules = rules

	hosts, err := parseHosts(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Hosts = hosts

	dnsCfg, err := parseDNS(rawCfg, hosts, rules)
	if err != nil {
		return nil, err
	}
	config.DNS = dnsCfg

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

func parseDNS(rawCfg *RawConfig, host *trie.DomainTrie[resolver.HostValue], rules []C.Rule) (*DNS, error) {
	cfg := rawCfg.DNS

	return &DNS{
		Enable: cfg.Enable,
	}, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, providersMap map[string]providerTypes.ProxyProvider, err error) {
	proxies = make(map[string]C.Proxy)
	providersMap = make(map[string]providerTypes.ProxyProvider)
	proxiesConfig := cfg.Proxy
	groupsConfig := cfg.ProxyGroup
	providersConfig := cfg.ProxyProvider

	var proxyList []string

	proxies["DIRECT"] = adapter.NewProxy(outbound.NewDirect())
	proxyList = append(proxyList, "DIRECT")

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

	for idx, mapping := range groupsConfig {
		groupName, existName := mapping["name"].(string)
		if !existName {
			return nil, nil, fmt.Errorf("代理组 %d: 名称缺失", idx)
		}
		proxyList = append(proxyList, groupName)
	}
	if err := proxyGroupDagSort(groupsConfig); err != nil {
		return nil, nil, err
	}

	for name := range providersConfig {
		if name == provider.ReservedName {
			return nil, nil, fmt.Errorf("不能定义provider: %s", provider.ReservedName)
		}
	}

	for idx, mapping := range groupsConfig {
		group, err := outboundgroup.ParseProxyGroup(mapping, proxies, providersMap)
		if err != nil {
			return nil, nil, fmt.Errorf("代理组[%d]: %w", idx, err)
		}

		groupName := group.Name()
		if _, exist := proxies[groupName]; exist {
			return nil, nil, fmt.Errorf("代理组 %s: 名称重复", groupName)
		}

		proxies[groupName] = adapter.NewProxy(group)
	}

	var ps []C.Proxy
	// proxyList 包括代理和代理组
	for _, v := range proxyList {
		if proxies[v].Type() == C.Pass {
			continue
		}

		ps = append(ps, proxies[v])
	}

	return proxies, providersMap, nil

}

func parseSubRules(cfg *RawConfig, proxies map[string]C.Proxy) (subRules map[string][]C.Rule, err error) {
	subRules = map[string][]C.Rule{}
	for name := range cfg.SubRules {
		subRules[name] = make([]C.Rule, 0)
	}

	for name, rawRules := range cfg.SubRules {
		if len(name) == 0 {
			return nil, fmt.Errorf("sub-rule 名称为空")
		}

		var rules []C.Rule
		rules, err = parseRules(rawRules, proxies, subRules, fmt.Sprintf("sub-rules[%s]", name))
		if err != nil {
			return nil, err
		}

		subRules[name] = rules
	}

	return
}

func parseRules(rulesConfig []string, proxies map[string]C.Proxy, subRules map[string][]C.Rule, format string) ([]C.Rule, error) {
	var rules []C.Rule
	for idx, line := range rulesConfig {
		rule := trimArr(strings.Split(line, ","))
		var (
			payload  string
			target   string
			params   []string
			ruleName = strings.ToUpper(rule[0])
		)

		l := len(rule)

		if l < 2 {
			return nil, fmt.Errorf("%s[%d] [%s]失败: 格式无效", format, idx, line)
		}

		if l < 4 {
			rule = append(rule, make([]string, 4-l)...)
		}

		if l >= 3 {
			l = 3
			payload = rule[1]
		}

		target = rule[l-1]
		params = rule[l:]

		params = trimArr(params)
		parsed, parseErr := R.ParseRule(ruleName, payload, target, params, subRules)
		if parseErr != nil {
			return nil, fmt.Errorf("%s[%d] [%s] 失败: %s", format, idx, line, parseErr.Error())
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}

func parseHosts(cfg *RawConfig) (*trie.DomainTrie[resolver.HostValue], error) {
	tree := trie.New[resolver.HostValue]()

	hostValue, _ := resolver.NewHostValueByIPs([]netip.Addr{netip.AddrFrom4([4]byte{127, 0, 0, 1})})
	if err := tree.Insert("localhost", hostValue); err != nil {
		log.Errorln("添加localhost到host失败: %s", err.Error())
	}

	if len(cfg.Hosts) != 0 {
		for range cfg.Hosts {

		}

	}
	tree.Optimize()

	return tree, nil
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{
		AllowLan:    false,
		BindAddress: "*",
		Proxy:       []map[string]any{},
		ProxyGroup:  []map[string]any{},
		Rule:        []string{},
	}

	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}

	return rawCfg, nil
}

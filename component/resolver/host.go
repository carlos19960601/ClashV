package resolver

import (
	"errors"
	"net/netip"
	"strings"

	"github.com/carlos19960601/ClashV/common/utils"
	"github.com/carlos19960601/ClashV/component/trie"
	"github.com/zhangyunhao116/fastrand"
)

type Hosts struct {
	*trie.DomainTrie[HostValue]
}

func NewHosts(hosts *trie.DomainTrie[HostValue]) Hosts {
	return Hosts{hosts}
}

func (h *Hosts) Search(domain string, isDomain bool) (*HostValue, bool) {
	if value := h.DomainTrie.Search(domain); value != nil {

	}

	return nil, false
}

type HostValue struct {
	IsDomain bool
	IPs      []netip.Addr
	Domain   string
}

func NewHostValue(value any) (HostValue, error) {
	isDomain := true
	ips := make([]netip.Addr, 0)
	domain := ""
	if valueArr, err := utils.ToStringSlice(value); err != nil {
		return HostValue{}, err
	} else {
		if len(valueArr) > 1 {
			isDomain = false
			for _, str := range valueArr {
				if ip, err := netip.ParseAddr(str); err == nil {
					ips = append(ips, ip)
				} else {
					return HostValue{}, nil
				}
			}
		} else if len(valueArr) == 1 {
			host := valueArr[0]
			if ip, err := netip.ParseAddr(host); err == nil {
				ips = append(ips, ip)
				isDomain = false
			} else {
				domain = host
			}
		}
	}

	if isDomain {
		return NewHostValueByDomain(domain)
	} else {
		return NewHostValueByIPs(ips)
	}

}

func NewHostValueByIPs(ips []netip.Addr) (HostValue, error) {
	if len(ips) == 0 {
		return HostValue{}, errors.New("ip列表为空")
	}

	return HostValue{
		IsDomain: false,
		IPs:      ips,
	}, nil
}

func NewHostValueByDomain(domain string) (HostValue, error) {
	domain = strings.Trim(domain, ".")
	item := strings.Split(domain, ".")
	if len(item) < 2 {
		return HostValue{}, errors.New("无效的域名")
	}

	return HostValue{
		IsDomain: true,
		Domain:   domain,
	}, nil
}

func (hv HostValue) RandIP() (netip.Addr, error) {
	if hv.IsDomain {
		return netip.Addr{}, errors.New("value type is error")
	}

	return hv.IPs[fastrand.Intn(len(hv.IPs))], nil
}

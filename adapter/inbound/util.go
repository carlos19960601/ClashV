package inbound

import (
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	C "github.com/carlos19960601/ClashV/constant"
)

func parseHTTPAddr(request *http.Request) *C.Metadata {
	host := request.URL.Hostname()
	port := request.URL.Port()

	if port == "" {
		port = "80"
	}

	host = strings.TrimRight(host, ".")

	var uint16Port uint16
	if port, err := strconv.ParseUint(port, 10, 16); err == nil {
		uint16Port = uint16(port)
	}

	metadata := &C.Metadata{
		NetWork: C.TCP,
		Host:    host,
		DstIP:   netip.Addr{},
		DstPort: uint16Port,
	}

	return metadata
}

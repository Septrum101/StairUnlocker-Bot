package provider

import (
	"fmt"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// hbo
type hbo struct {
}

func (h *hbo) create() AbsStream {
	return new(hbo)
}

func (h *hbo) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "HBO"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.hbomax.com")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if !strings.Contains(resp.RawResponse.Request.URL.Path, "geo") {
		s.Unlock = true
	}
	return
}

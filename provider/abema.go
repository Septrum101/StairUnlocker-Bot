package provider

import (
	"fmt"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// abema
type abema struct {
}

func (a *abema) create() AbsStream {
	return new(abema)
}

func (a *abema) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "Abema"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://api.abema.io/v1/ip/check?device=android")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if strings.Contains(resp.String(), "Country") {
		s.Unlock = true
	}
	return
}

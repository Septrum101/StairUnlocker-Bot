package provider

import (
	"fmt"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// disney
type disney struct {
}

func (d *disney) create() AbsStream {
	return new(disney)
}

func (d *disney) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "Disney Plus"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.disneyplus.com")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Unlock = true
	}
	return
}

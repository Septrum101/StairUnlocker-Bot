package provider

import (
	"fmt"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// connectivity
type connectivity struct {
}

func (c *connectivity) create() AbsStream {
	return new(connectivity)
}

func (c *connectivity) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "<Connectivity>"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "http://www.gstatic.com/generate_204")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Unlock = true
	}
	return
}

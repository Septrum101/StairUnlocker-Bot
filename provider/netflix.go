package provider

import (
	"fmt"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// netflix
type netflix struct {
}

func (n *netflix) create() AbsStream {
	return new(netflix)
}

func (n *netflix) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "Netflix"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.netflix.com/title/70143836")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Unlock = true
	}
	return
}

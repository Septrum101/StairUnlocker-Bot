package provider

import (
	"fmt"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// youtube
type youtube struct {
}

func (y *youtube) create() AbsStream {
	return new(youtube)
}

func (y *youtube) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "Youtube Premium"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://music.youtube.com")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Unlock = true
	}
	return
}

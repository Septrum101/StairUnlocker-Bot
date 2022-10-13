package provider

import (
	"encoding/json"
	"fmt"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// tvb
type tvb struct {
}

func (t *tvb) create() AbsStream {
	return new(tvb)
}

func (t *tvb) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "TVB"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.mytvsuper.com/api/auth/getSession/self")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	r := make(map[string]int)
	json.Unmarshal(resp.Body(), &r)

	if r["region"] == 1 {
		s.Unlock = true
	}
	return
}

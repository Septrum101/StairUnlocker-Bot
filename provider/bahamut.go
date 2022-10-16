package provider

import (
	"fmt"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

// bahamut
type bahamut struct {
}

func (b *bahamut) create() AbsStream {
	return new(bahamut)
}

func (b *bahamut) IsUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "Bahamut"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667")
	s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	if err != nil {
		return
	}
	if strings.Contains(resp.String(), "animeSn") {
		s.Unlock = true
	}
	return
}

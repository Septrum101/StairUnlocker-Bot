package utils

import (
	"fmt"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
)

type absStream interface {
	isUnlock(p *C.Proxy) (StreamData, error)
}

type unlockFactory interface {
	create() absStream
}

type netflix struct {
}

func (n *netflix) create() absStream {
	return new(netflix)
}

func (n *netflix) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "Netflix"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.netflix.com/title/70143836")
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

type hbo struct {
}

func (h *hbo) create() absStream {
	return new(hbo)
}

func (h *hbo) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "HBO"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.hbomax.com")
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

type youtube struct {
}

func (y *youtube) create() absStream {
	return new(youtube)
}

func (y *youtube) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "Youtube Premium"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://music.youtube.com")
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

type disney struct {
}

func (d *disney) create() absStream {
	return new(disney)
}

func (d *disney) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "Disney Plus"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.disneyplus.com")
	if err != nil {
		return
	}
	if resp.StatusCode() < 300 {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

type tvb struct {
}

func (t *tvb) create() absStream {
	return new(tvb)
}

func (t *tvb) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "TVB"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://www.mytvsuper.com/iptest.php")
	if err != nil {
		return
	}
	if strings.Contains(resp.String(), "HK") {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

type abema struct {
}

func (a *abema) create() absStream {
	return new(abema)
}

func (a *abema) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "Abema"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://api.abema.io/v1/ip/check?device=android")
	if err != nil {
		return
	}
	if strings.Contains(resp.String(), "Country") {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

type bahamut struct {
}

func (b *bahamut) create() absStream {
	return new(bahamut)
}

func (b *bahamut) isUnlock(p *C.Proxy) (s StreamData, err error) {
	s.Name = "bahamut"
	s.ProxyName = (*p).Name()
	start := time.Now()
	resp, err := getURLResp(p, "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667")
	if err != nil {
		return
	}
	if strings.Contains(resp.String(), "animeSn") {
		s.Latency = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
		s.Unlock = true
	}
	return
}

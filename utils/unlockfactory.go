package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

type absStream interface {
	isUnlock(p *C.Proxy) (model.StreamData, error)
}

type unlockFactory interface {
	create() absStream
}

func NewStreamList() []absStream {
	return []absStream{
		unlockFactory(new(netflix)).create(),
		unlockFactory(new(hbo)).create(),
		unlockFactory(new(youtube)).create(),
		unlockFactory(new(disney)).create(),
		unlockFactory(new(tvb)).create(),
		unlockFactory(new(abema)).create(),
		unlockFactory(new(bahamut)).create(),
	}
}

// netflix
type netflix struct {
}

func (n *netflix) create() absStream {
	return new(netflix)
}

func (n *netflix) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
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

// youtube
type youtube struct {
}

func (y *youtube) create() absStream {
	return new(youtube)
}

func (y *youtube) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
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

// disney
type disney struct {
}

func (d *disney) create() absStream {
	return new(disney)
}

func (d *disney) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
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

// hbo
type hbo struct {
}

func (h *hbo) create() absStream {
	return new(hbo)
}

func (h *hbo) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
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

// tvb
type tvb struct {
}

func (t *tvb) create() absStream {
	return new(tvb)
}

func (t *tvb) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
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

//abema
type abema struct {
}

func (a *abema) create() absStream {
	return new(abema)
}

func (a *abema) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
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

// bahamut
type bahamut struct {
}

func (b *bahamut) create() absStream {
	return new(bahamut)
}

func (b *bahamut) isUnlock(p *C.Proxy) (s model.StreamData, err error) {
	s.Name = "bahamut"
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

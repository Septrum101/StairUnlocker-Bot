package user

import (
	"errors"
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/thank243/StairUnlocker-Bot/config"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (u *User) generateProxies(apiURL string) (proxies map[string]C.Proxy, unmarshalProxies *config.RawConfig, err error) {
	log.Infoln("[ID: %d] Converting from API server.", u.ID)
	_, err = u.SendMessage("Converting from API server.", true)
	pList, err := u.convertAPI(apiURL)
	if err != nil {
		u.Data.SubURL = ""
		return
	}
	u.Data.LastCheck = time.Now().Unix()
	unmarshalProxies, err = config.UnmarshalRawConfig(pList)
	if err != nil {
		return
	}
	//proxiesTest(u)
	// compatible clash-core 1.9.0
	for i := range unmarshalProxies.Proxy {
		for k := range unmarshalProxies.Proxy[i] {
			switch k {
			case "ws-path":
				unmarshalProxies.Proxy[i]["ws-opts"] = map[string]interface{}{"path": unmarshalProxies.Proxy[i]["ws-path"]}
				delete(unmarshalProxies.Proxy[i], "ws-path")
			case "ws-header":
				unmarshalProxies.Proxy[i]["ws-opts"] = map[string]interface{}{"ws-header": unmarshalProxies.Proxy[i]["ws-header"]}
				delete(unmarshalProxies.Proxy[i], "ws-header")
			}
		}
	}
	proxies, err = u.parseProxies(unmarshalProxies)
	return
}

func (u *User) convertAPI(apiURL string) (re []byte, err error) {
	baseUrl, err := url.Parse(apiURL)
	baseUrl.Path += "sub"
	params := url.Values{}
	params.Add("target", "clash")
	params.Add("list", strconv.FormatBool(true))
	params.Add("emoji", strconv.FormatBool(false))
	params.Add("url", u.Data.SubURL)
	baseUrl.RawQuery = params.Encode()
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, baseUrl.String(), nil)
	req.Header.Set("User-Agent", "ClashforWindows/0.19.2")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	re, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Errorln("[ID: %d] %s", u.ID, re)
		err = errors.New(string(re))
		return
	}
	return
}

func (u *User) parseProxies(cfg *config.RawConfig) (proxies map[string]C.Proxy, err error) {
	if cfg == nil {
		err = fmt.Errorf("the original converted URL must be used for clash")
		return
	}
	proxies = make(map[string]C.Proxy)
	proxiesConfig := cfg.Proxy
	for idx, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			return nil, fmt.Errorf("proxy %d: %w", idx, err)
		}
		if _, exist := proxies[proxy.Name()]; exist {
			return nil, fmt.Errorf("proxy %s is the duplicate name", proxy.Name())
		}
		proxies[proxy.Name()] = proxy
	}
	return
}

package user

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/thank243/StairUnlocker-Bot/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (u *User) generateProxies(apiURL string) (proxies map[string]C.Proxy, unmarshalProxies *config.RawConfig, err error) {
	log.Infoln("[ID: %d] Converting from API server.", u.ID)
	_, err = u.Send("Converting from API server.", true)
	pList, err := u.convertAPI(apiURL)
	if err != nil {
		return
	}
	unmarshalProxies, _ = config.UnmarshalRawConfig(pList)
	proxies, err = u.parseProxies(unmarshalProxies)
	return
}

func (u *User) convertAPI(apiURL string) (re []byte, err error) {
	baseUrl, err := url.Parse(apiURL)
	baseUrl.Path += "sub"
	params := url.Values{}
	params.Add("target", "clash")
	params.Add("list", strconv.FormatBool(true))
	params.Add("url", u.Data.SubURL)
	params.Add("emoji", strconv.FormatBool(false))
	baseUrl.RawQuery = params.Encode()
	reqs, err := http.Get(baseUrl.String())
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(reqs.Body)
	re, _ = ioutil.ReadAll(reqs.Body)
	if reqs.StatusCode != 200 {
		log.Errorln("[ID: %d] %s", u.ID, re)
		err = fmt.Errorf(string(re))
		return
	}
	return
}

func (u *User) parseProxies(cfg *config.RawConfig) (proxies map[string]C.Proxy, err error) {
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

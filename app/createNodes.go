package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/go-resty/resty/v2"

	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/model"
)

func (u *User) buildProxies(subUrl string) (proxies map[string]C.Proxy, err error) {
	log.Infoln("[ID: %d] Converting from API server.", u.ID)
	_, err = u.SendMessage("Converting from API server.")
	pList, err := u.convertAPI(subUrl)
	if err != nil {
		u.Data.SubURL = ""
		return
	}
	unmarshalProxies, err := config.UnmarshalRawConfig(pList)
	if err != nil {
		return
	}
	if len(unmarshalProxies.Proxy) == 0 {
		log.Errorln("[ID: %d] %s", u.ID, "No nodes were found!")
		err = errors.New("no nodes were found")
		return
	}
	if len(unmarshalProxies.Proxy) > 1024 {
		log.Errorln("[ID: %d] %s", u.ID, "Too many nodes at the same time, Please reduce nodes less than 1024.")
		err = errors.New("too many nodes")
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

func (u *User) convertAPI(subUrl string) (re []byte, err error) {
	if subUrl == "" {
		subUrl = u.Data.SubURL
	}
	resp, err := resty.New().SetHeader("User-Agent", "ClashforWindows/0.19.6").SetRetryCount(3).
		SetQueryParams(map[string]string{
			"target":      "clash",
			"append_type": strconv.FormatBool(true),
			"list":        strconv.FormatBool(true),
			"emoji":       strconv.FormatBool(false),
			"url":         subUrl},
		).R().Get(fmt.Sprintf("%s/sub", model.BotCfg.ConverterAPI))

	re = resp.Body()
	if resp.StatusCode() != 200 {
		log.Errorln("[ID: %d] %s", u.ID, re)
		err = errors.New(string(re))
		return
	}
	return
}

func (u *User) parseProxies(cfg *model.RawConfig) (proxies map[string]C.Proxy, err error) {
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

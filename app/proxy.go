package app

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/model"
)

func (u *User) buildProxies(subUrl string) (proxies map[string]C.Proxy, err error) {
	log.Infoln("[ID: %d] Converting from API server.", u.ID)

	resp, err := convertAPI(subUrl)
	if err != nil {
		log.Errorln("[ID: %d] %s", u.ID, err)
		return
	}

	rawConfig, err := config.UnmarshalRawConfig(resp)
	if err != nil {
		return
	}
	if len(rawConfig.Proxy) == 0 {
		log.Errorln("[ID: %d] %s", u.ID, "No nodes were found!")
		err = fmt.Errorf("no nodes were found")
		return
	}
	if len(rawConfig.Proxy) > 1024 {
		log.Errorln("[ID: %d] %s", u.ID, "Too many nodes at the same time, Please reduce nodes less than 1024.")
		err = fmt.Errorf("too many nodes")
		return
	}
	// proxiesTest(u)
	// compatible clash-core 1.9.0
	for i := range rawConfig.Proxy {
		for k := range rawConfig.Proxy[i] {
			switch k {
			case "ws-path":
				rawConfig.Proxy[i]["ws-opts"] = map[string]interface{}{"path": rawConfig.Proxy[i]["ws-path"]}
				delete(rawConfig.Proxy[i], "ws-path")
			case "ws-header":
				rawConfig.Proxy[i]["ws-opts"] = map[string]interface{}{"ws-header": rawConfig.Proxy[i]["ws-header"]}
				delete(rawConfig.Proxy[i], "ws-header")
			}
		}
	}
	proxies, err = parseProxies(rawConfig)
	return
}

func convertAPI(subUrl string) ([]byte, error) {
	unescapeUrl, _ := url.QueryUnescape(subUrl)
	resp, err := resty.New().SetHeader("User-Agent", "ClashforWindows/0.19.6").SetRetryCount(3).
		SetQueryParams(map[string]string{
			"target":      "clash",
			"append_type": strconv.FormatBool(true),
			"list":        strconv.FormatBool(true),
			"emoji":       strconv.FormatBool(false),
			"url":         unescapeUrl,
		}).R().Get(fmt.Sprintf("%s/sub", viper.GetString("converterAPI")))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() > 299 {
		return nil, fmt.Errorf(resp.String())
	}
	return resp.Body(), nil
}

func parseProxies(cfg *model.RawConfig) (proxies map[string]C.Proxy, err error) {
	proxies = make(map[string]C.Proxy)
	for idx, mapping := range cfg.Proxy {
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

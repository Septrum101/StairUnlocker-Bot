package config

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies,flow"`
}

type SuConfig struct {
	ConverterAPI string       `yaml:"converterAPI"`
	MaxConn      int          `yaml:"maxConn"`
	MaxOnline    int          `yaml:"maxOnline"`
	LogLevel     log.LogLevel `yaml:"log_level"`
	Telegram     *telegram
}
type telegram struct {
	TelegramToken string `yaml:"telegramToken,omitempty"`
	ChatID        int64  `yaml:"chatID,omitempty"`
	Secret        string `yaml:"secret,omitempty"`
}

func Init(cfgPath *string) (s *SuConfig) {
	//initial config.yaml
	var buf []byte
	if *cfgPath != "" {
		buf, _ = ioutil.ReadFile(*cfgPath)
	} else {
		_, err := os.Stat("config.yaml")
		if err != nil {
			b, _ := ioutil.ReadFile("config.example.yaml")
			_ = ioutil.WriteFile("config.yaml", b, 644)
		}
		buf, _ = ioutil.ReadFile("config.yaml")
	}
	var cfg SuConfig
	_ = yaml.Unmarshal(buf, &cfg)
	return &cfg
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{}
	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}
	return rawCfg, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, err error) {
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
	return proxies, err
}

func GenerateProxies(apiURL, subURL string) (proxies map[string]C.Proxy, unmarshalProxies *RawConfig, err error) {
	log.Infoln("Converting from API server.")
	pList, err := convertAPI(apiURL, subURL)
	if err != nil {
		return
	}
	unmarshalProxies, _ = UnmarshalRawConfig(pList)
	proxies, err = parseProxies(unmarshalProxies)
	return proxies, unmarshalProxies, err
}

func convertAPI(apiURL, subURL string) (p []byte, err error)  {
	baseUrl, err := url.Parse(apiURL)
	baseUrl.Path += "sub"
	params := url.Values{}
	params.Add("target", "clash")
	params.Add("list", "true")
	params.Add("url", subURL)
	baseUrl.RawQuery = params.Encode()
	reqs, err := http.Get(baseUrl.String())
	if err != nil {
		log.Errorln(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(reqs.Body)
	p, _ = ioutil.ReadAll(reqs.Body)
	if strings.Contains(string(p), "The following link doesn't contain any valid node info") {
		log.Errorln("The following link doesn't contain any valid node info.")
		err = fmt.Errorf("invalid link")
		return nil, err
	}
	return
}

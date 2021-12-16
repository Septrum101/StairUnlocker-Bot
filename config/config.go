package config

import (
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies,flow"`
}

type SuConfig struct {
	ConverterAPI  string       `yaml:"converterAPI"`
	MaxConn       int          `yaml:"maxConn"`
	MaxOnline     int          `yaml:"maxOnline"`
	LogLevel      log.LogLevel `yaml:"log_level"`
	Internal      int          `yaml:"internal"`
	TelegramToken string       `yaml:"telegramToken,omitempty"`
}

func Init(cfgPath *string) (cfg *SuConfig) {
	//initial config.yaml
	var (
		buf []byte
	)
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
	_ = yaml.Unmarshal(buf, &cfg)
	return
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{}
	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}
	return rawCfg, nil
}

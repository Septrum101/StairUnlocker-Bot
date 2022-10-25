package config

import (
	"github.com/Dreamacro/clash/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/thank243/StairUnlocker-Bot/model"
)

func init() {
	// initial config.yaml
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("fatal error config file: %w", err)
	}
}

func UnmarshalRawConfig(buf []byte) (*model.RawConfig, error) {
	rawConf := &model.RawConfig{}
	if err := yaml.Unmarshal(buf, rawConf); err != nil {
		return nil, err
	}
	return rawConf, nil
}

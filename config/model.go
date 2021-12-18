package config

import (
	"github.com/Dreamacro/clash/log"
)

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies,flow"`
}

type BotConfig struct {
	ConverterAPI  string       `yaml:"converterAPI"`
	MaxConn       int          `yaml:"maxConn"`
	MaxOnline     int          `yaml:"maxOnline"`
	LogLevel      log.LogLevel `yaml:"log_level"`
	Internal      int          `yaml:"internal"`
	TelegramToken string       `yaml:"telegramToken,omitempty"`
}

var (
	Version    bool
	Help       bool
	BotCfg     *BotConfig
	configPath string
)

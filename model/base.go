package model

type StreamData struct {
	Name      string
	ProxyName string
	Latency   string
	Unlock    bool
}

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies"`
}

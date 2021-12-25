package user

import "github.com/thank243/StairUnlocker-Bot/config"

func FilterProxies(netflixList []string, cfg *config.RawConfig) (netflixCfg config.RawConfig) {
	for idx := range netflixList {
		for i := range cfg.Proxy {
			if netflixList[idx] == cfg.Proxy[i]["name"] {
				netflixCfg.Proxy = append(netflixCfg.Proxy, cfg.Proxy[i])
				//删除已添加元素
				cfg.Proxy = append(cfg.Proxy[:i], cfg.Proxy[i+1:]...)
				break
			}
		}
	}
	return
}

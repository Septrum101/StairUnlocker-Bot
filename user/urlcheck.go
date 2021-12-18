package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

func (u *User) URLCheck() {
	var proxiesList []C.Proxy
	u.IsCheck = true
	proxies, unmarshalProxies, err := u.generateProxies(config.BotCfg.ConverterAPI)
	if err != nil {
		_ = u.Send(err.Error())
		return
	}
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	//同时连接数
	connNum := config.BotCfg.MaxConn
	if i := len(proxiesList); i < connNum {
		connNum = i
	}
	start := time.Now()
	netflixList, latency := utils.BatchCheck(proxiesList, connNum)
	//proxiesTest(netflixList,u)
	report := fmt.Sprintf("Total %d nodes, %d unlock nodes.\nElapsed time: %s", len(proxiesList), len(netflixList), time.Since(start).Round(time.Millisecond))
	log.Warnln(report)
	telegramReport := fmt.Sprintf("%s\nTimestamp: %s\n%s\n%s", report, time.Now().Round(time.Millisecond), strings.Repeat("-", 35), strings.Join(latency, "\n"))
	// todo upload file
	_, _ = yaml.Marshal(NetflixFilter(netflixList, unmarshalProxies))
	u.Data.CheckInfo = telegramReport
	_ = u.Send(telegramReport)
	u.IsCheck = false
}

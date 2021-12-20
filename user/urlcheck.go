package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

func (u *User) URLCheck() {
	var proxiesList []C.Proxy
	u.IsCheck = true
	defer func() { u.IsCheck = false }()
	proxies, unmarshalProxies, err := u.generateProxies(config.BotCfg.ConverterAPI)
	if err != nil {
		_ = u.EditMessage(u.MessageID, err.Error())
		return
	}
	_ = u.EditMessage(u.MessageID, "Checking nodes unlock status...")
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
	u.Data.CheckInfo = telegramReport
	_ = u.EditMessage(u.MessageID, telegramReport)
	// send proxies.yaml
	marshal, _ := yaml.Marshal(NetflixFilter(netflixList, unmarshalProxies))
	_, err = u.Bot.Send(tgbotapi.NewDocument(u.ID, tgbotapi.FileBytes{
		Name:  "proxies.yaml",
		Bytes: marshal,
	}))
	if err != nil {
		log.Errorln(err.Error())
	}
}

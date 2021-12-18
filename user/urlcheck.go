package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

type User struct {
	ID      int64
	Data    Data
	Bot     *tgBot.BotAPI
	IsCheck bool
}

type Data struct {
	LastCheck int64
	SubURL    string
	CheckInfo string
}

func (u *User) URLCheck(apiURL string, maxConn int) {
	var proxiesList []C.Proxy
	u.IsCheck = true
	proxies, unmarshalProxies, err := u.generateProxies(apiURL)
	if err != nil {
		_ = u.Send(err.Error())
		return
	}
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	//同时连接数
	connNum := maxConn
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

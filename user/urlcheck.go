package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/exporter"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"sort"
	"time"
)

func (u *User) URLCheck() {
	var proxiesList []C.Proxy
	u.IsCheck = true
	defer func() { u.IsCheck = false }()
	proxies, _, err := u.generateProxies(config.BotCfg.ConverterAPI)
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
	// 有效节点才开始测试
	if len(proxiesList) > 0 {
		start := time.Now()
		streamMediaList, latencyMap := utils.BatchCheck(proxiesList, connNum)
		sort.Strings(streamMediaList)

		//proxiesTest(streamMediaList,u)
		log.Warnln(fmt.Sprintf("Total %d nodesd. Elapsed time: %s", len(proxiesList), time.Since(start).Round(time.Millisecond)))
		report := fmt.Sprintf("Total %d nodes.\nElapsed time: %s", len(proxiesList), time.Since(start).Round(time.Millisecond))
		telegramReport := fmt.Sprintf("%s\nTimestamp: %s\n", report, time.Now().Round(time.Millisecond))
		u.Data.CheckInfo = telegramReport
		// update message, not send a new message
		_ = u.EditMessage(u.MessageID, telegramReport)
		// send result image
		buffer, err := exporter.Export(latencyMap)
		if err != nil {
			return
		}
		_, err = u.Bot.Send(tgbotapi.NewDocument(u.ID, tgbotapi.FileBytes{
			Name:  fmt.Sprintf("stairunlocker_bot_%d.png", time.Now().Unix()),
			Bytes: buffer.Bytes(),
		}))
	}
}

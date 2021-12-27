package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"strings"
	"time"
)

func statistic(streamMediaMap *map[string][]uint16) map[string]int {
	statMap := make(map[string]int)
	for i := range *streamMediaMap {
		for idx := range (*streamMediaMap)[i] {
			switch idx {
			case 0:
				if (*streamMediaMap)[i][idx] != 0 {
					statMap["Netflix"]++
				}
			case 1:
				if (*streamMediaMap)[i][idx] != 0 {
					statMap["HBO"]++
				}
			case 2:
				if (*streamMediaMap)[i][idx] != 0 {
					statMap["Disney Plus"]++
				}
			case 3:
				if (*streamMediaMap)[i][idx] != 0 {
					statMap["Youtube Premium"]++
				}
			}
		}
	}
	return statMap
}

func (u *User) URLCheck() {
	var proxiesList []C.Proxy
	u.IsCheck = true
	defer func() {
		u.IsCheck = false
	}()
	proxies, _, err := u.generateProxies(config.BotCfg.ConverterAPI)
	if err != nil {
		_ = u.EditMessage(u.MessageID, err.Error())
		return
	}
	// animation while waiting test.
	go func() {
		log.Infoln("[ID: %d]: Checking nodes unlock status.", u.ID)
		count := 0
		for u.IsCheck {
			count++
			if count > 5 {
				count = 0
			}
			_ = u.EditMessage(u.MessageID, fmt.Sprintf("Checking nodes unlock status%s", strings.Repeat(".", count)))
			time.Sleep(500 * time.Millisecond)
		}
		return
	}()

	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	// 同时连接数
	connNum := config.BotCfg.MaxConn
	if i := len(proxiesList); i < connNum {
		connNum = i
	}
	// 必须包含节点
	if len(proxiesList) > 0 {
		start := time.Now()
		streamMediaUnlockMap := utils.BatchCheck(proxiesList, connNum)
		u.IsCheck = false
		report := fmt.Sprintf("Total %d nodes tested\nElapsed time: %s", len(proxiesList), time.Since(start).Round(time.Millisecond))
		// save test results.
		result := statistic(&streamMediaUnlockMap)
		var finalStr string
		for k, v := range result {
			finalStr += fmt.Sprintf("%s: %d\n", k, v)
		}
		telegramReport := fmt.Sprintf("StairUnlocker Bot Bulletin:\n%s\n%sTimestamp: %s\n%s", report, finalStr, time.Now().Round(time.Second), strings.Repeat("-", 30))
		u.Data.CheckInfo = telegramReport
		log.Warnln(report)
		_ = u.EditMessage(u.MessageID, "Uploading PNG file...")
		buffer, err := generatePNG(streamMediaUnlockMap)
		if err != nil {
			return
		}
		// send result image
		docToSend := tgbotapi.NewDocument(u.ID, tgbotapi.FileBytes{
			Name:  fmt.Sprintf("stairunlocker_bot_result_%d.png", time.Now().Unix()),
			Bytes: buffer.Bytes(),
		})
		docToSend.Caption = telegramReport + "\n@stairunlock_test_bot\nProject: https://git.io/Jyl5l"
		_, err = u.Bot.Send(docToSend)
		_ = u.DeleteMessage(u.MessageID)
		//proxiesTest(u)
	}
}

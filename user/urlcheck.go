package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"sort"
	"strings"
	"time"
)

func sendUserCheckStatus(u *User, c chan bool) {
	log.Infoln("[ID: %d] Checking nodes unlock status.", u.ID)
	count := 0
	for {
		select {
		case <-c:
			return
		default:
			count++
			if count > 5 {
				count = 0
			}
			_ = u.EditMessage(u.MessageID, fmt.Sprintf("Checking nodes unlock status%s", strings.Repeat(".", count)))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func statistic(streamMediaList *[]utils.CheckData) map[string]int {
	statMap := make(map[string]int)
	statMap["Netflix"], statMap["HBO"], statMap["Disney Plus"], statMap["Youtube Premium"] = 0, 0, 0, 0
	for i := range *streamMediaList {
		switch (*streamMediaList)[i].StreamMedia {
		case "Netflix":
			statMap["Netflix"]++
		case "HBO":
			statMap["HBO"]++
		case "Disney Plus":
			statMap["Disney Plus"]++
		case "Youtube Premium":
			statMap["Youtube Premium"]++
		}
	}
	return statMap
}

func (u *User) URLCheck() {
	var proxiesList []C.Proxy
	c := make(chan bool)
	u.IsCheck = true
	defer func() {
		u.IsCheck = false
		close(c)
	}()

	proxies, _, err := u.generateProxies(config.BotCfg.ConverterAPI)
	if err != nil {
		_ = u.EditMessage(u.MessageID, err.Error())
		return
	}
	// animation while waiting test.
	go sendUserCheckStatus(u, c)

	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	//	Max Connection at the same time.
	connNum := config.BotCfg.MaxConn
	if i := len(proxiesList); i < connNum {
		connNum = i
	}
	// Must have valid node.
	if len(proxiesList) > 0 {
		start := time.Now()
		streamMediaUnlockList := utils.BatchCheck(proxiesList, connNum)
		c <- true
		report := fmt.Sprintf("Total %d nodes tested\nElapsed time: %s", len(proxiesList), time.Since(start).Round(time.Millisecond))
		// save test results.
		var finalStr string
		i := 0
		var nameList []string
		statisticMap := statistic(&streamMediaUnlockList)
		for k := range statisticMap {
			nameList = append(nameList, k)
			i++
		}
		sort.Strings(nameList)
		for i := range nameList {
			finalStr += fmt.Sprintf("%s: %d\n", nameList[i], statisticMap[nameList[i]])
		}
		telegramReport := fmt.Sprintf("StairUnlocker Bot Bulletin:\n%s\n%sTimestamp: %s\n%s", report, finalStr, time.Now().Round(time.Second), strings.Repeat("-", 30))
		u.Data.CheckInfo = telegramReport
		log.Warnln("[ID: %d] %s", u.ID, report)
		_ = u.EditMessage(u.MessageID, "Uploading PNG file...")
		streamMediaUnlockMap := make(map[string][]string)
		for i := range proxiesList {
			streamMediaUnlockMap[proxiesList[i].Name()] = make([]string, 4)
		}
		for idx := range streamMediaUnlockList {
			switch streamMediaUnlockList[idx].StreamMedia {
			case "Netflix":
				streamMediaUnlockMap[streamMediaUnlockList[idx].ProxyName][0] = streamMediaUnlockList[idx].Latency
			case "HBO":
				streamMediaUnlockMap[streamMediaUnlockList[idx].ProxyName][1] = streamMediaUnlockList[idx].Latency
			case "Disney Plus":
				streamMediaUnlockMap[streamMediaUnlockList[idx].ProxyName][2] = streamMediaUnlockList[idx].Latency
			case "Youtube Premium":
				streamMediaUnlockMap[streamMediaUnlockList[idx].ProxyName][3] = streamMediaUnlockList[idx].Latency
			}
		}
		buffer, err := generatePNG(streamMediaUnlockMap)
		if err != nil {
			return
		}
		// send result image
		wrapPNG := tgbotapi.NewDocument(u.ID, tgbotapi.FileBytes{
			Name:  fmt.Sprintf("stairunlocker_bot_result_%d.png", time.Now().Unix()),
			Bytes: buffer.Bytes(),
		})
		wrapPNG.Caption = fmt.Sprintf("%s\n@stairunlock_test_bot\nProject: https://git.io/Jyl5l", telegramReport)
		_, err = u.Bot.Send(wrapPNG)
		_ = u.DeleteMessage(u.MessageID)
		//proxiesTest(u)
	}
}

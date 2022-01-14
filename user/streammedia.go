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

func statistic(streamMediaList *[]utils.CheckData) map[string]int {
	statMap := make(map[string]int)
	// initial 0 for each stream media
	for i := range utils.GetCheckParams() {
		statMap[utils.GetCheckParams()[i].CheckName] = 0
	}
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
		case "TVB":
			statMap["TVB"]++
		case "Abema":
			statMap["Abema"]++
		case "Bahamut":
			statMap["Bahamut"]++
		}
	}
	return statMap
}

func (u *User) StreamMedia() {
	u.IsCheck = true
	var proxiesList []C.Proxy
	checkFlag := make(chan bool)
	defer func() {
		u.Data.LastCheck = time.Now().Unix()
		u.IsCheck = false
		close(checkFlag)
	}()

	proxies, _, err := u.generateProxies(config.BotCfg.ConverterAPI)
	if err != nil {
		_ = u.EditMessage(u.MessageID, err.Error())
		return
	}
	// animation while waiting test.
	go u.statusMessage("Checking nodes unlock status", checkFlag)

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
		unlockList := utils.BatchCheck(proxiesList, connNum)
		checkFlag <- true
		report := fmt.Sprintf("Total %d nodes, Duration: %s", len(proxiesList), time.Since(start).Round(time.Millisecond))

		var nameList []string
		i := 0
		statisticMap := statistic(&unlockList)
		for k := range statisticMap {
			nameList = append(nameList, k)
			i++
		}
		sort.Strings(nameList)
		var finalStr string
		for i := range nameList {
			finalStr += fmt.Sprintf("%s: %d\n", nameList[i], statisticMap[nameList[i]])
		}
		telegramReport := fmt.Sprintf("StairUnlocker Bot %s Bulletin:\n%s\n%sTimestamp: %s\n%s", C.Version, report, finalStr, time.Now().UTC().Format(time.RFC3339), strings.Repeat("-", 25))
		// save test results.
		u.Data.CheckInfo = telegramReport
		log.Warnln("[ID: %d] %s", u.ID, report)
		_ = u.EditMessage(u.MessageID, "Uploading PNG file...")
		unlockMap := make(map[string][]string)
		for i := range proxiesList {
			unlockMap[proxiesList[i].Name()] = make([]string, 7)
		}
		for idx := range unlockList {
			switch unlockList[idx].StreamMedia {
			case "Netflix":
				unlockMap[unlockList[idx].ProxyName][0] = unlockList[idx].Latency
			case "HBO":
				unlockMap[unlockList[idx].ProxyName][1] = unlockList[idx].Latency
			case "Disney Plus":
				unlockMap[unlockList[idx].ProxyName][2] = unlockList[idx].Latency
			case "Youtube Premium":
				unlockMap[unlockList[idx].ProxyName][3] = unlockList[idx].Latency
			case "TVB":
				unlockMap[unlockList[idx].ProxyName][4] = unlockList[idx].Latency
			case "Abema":
				unlockMap[unlockList[idx].ProxyName][5] = unlockList[idx].Latency
			case "Bahamut":
				unlockMap[unlockList[idx].ProxyName][6] = unlockList[idx].Latency
			}
		}

		buffer, err := generatePNG(unlockMap)
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
	}
}

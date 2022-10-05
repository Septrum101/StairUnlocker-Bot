package user

import (
	"fmt"
	"sort"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

func statistic(streamMediaList *[]utils.StreamData) map[string]int {
	statMap := make(map[string]int)
	for i := range *streamMediaList {
		statMap[(*streamMediaList)[i].Name]++
		if !(*streamMediaList)[i].Unlock {
			statMap[(*streamMediaList)[i].Name]--
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
		statisticMap := statistic(&unlockList)
		for k := range statisticMap {
			nameList = append(nameList, k)
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

		buffer, err := generatePNG(unlockList, nameList)
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

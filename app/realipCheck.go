package app

import (
	"fmt"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/thank243/StairUnlocker-Bot/model"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

func (u *User) realIP(subUrl string) error {
	u.IsCheck = true
	checkFlag := make(chan bool)
	defer func() {
		u.Data.LastCheck = time.Now().Unix()
		u.IsCheck = false
		close(checkFlag)
	}()

	proxies, err := u.buildProxies(subUrl)
	if err != nil {
		u.EditMessage(u.editMsgID, err.Error())
		return err
	}
	if subUrl != "" {
		u.Data.SubURL = subUrl
	}

	var proxiesList []C.Proxy
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	// animation status
	go u.loading("Retrieving IP information", checkFlag)

	start := time.Now()
	inbound, outbound := utils.GetIPList(proxiesList, model.BotCfg.MaxConn)
	checkFlag <- true
	log.Warnln("[ID: %d] Total %d nodes: inbounds: %d -> outbounds: %d", u.ID, len(proxies), len(inbound), len(outbound))
	ipStatTitle := fmt.Sprintf("StairUnlocker Bot %s Bulletin:\nTotal %d nodes, Duration: %s\ninbound IP: %d\noutbound IP: %d\nTimestamp: %s", C.Version, len(proxies), time.Since(start).Round(time.Millisecond), len(inbound), len(outbound), time.Now().UTC().Format(time.RFC3339))
	ipStat := fmt.Sprintf("StairUnlocker Bot %s Bulletin:\nEntrypoint IP: ", C.Version)
	for _, v := range inbound {
		ipStat += "\n" + v
	}
	ipStat += "\n\nEndpoint IP: "
	for _, v := range outbound {
		ipStat += "\n" + v
	}
	warpFile := tgBot.NewDocument(u.ID, tgBot.FileBytes{
		Name:  fmt.Sprintf("stairunlocker_bot_realIP_%d.txt", time.Now().Unix()),
		Bytes: []byte(ipStat),
	})
	warpFile.Caption = fmt.Sprintf("%s\n%s\n@stairunlock_test_bot\nProject: https://git.io/Jyl5l", ipStatTitle, strings.Repeat("-", 25))
	u.s.Bot.Send(warpFile)
	u.DeleteMessage(u.editMsgID)
	return nil
}

package user

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"net/url"
	"strings"
	"time"
)

func (u *User) RealIP(messageText string) {
	u.IsCheck = true
	checkFlag := make(chan bool)
	defer func() {
		u.IsCheck = false
		close(checkFlag)
	}()
	var subURL *url.URL
	trimStr := strings.TrimSpace(strings.ReplaceAll(messageText, "/ip", ""))
	subURL, err := url.Parse(trimStr)
	if err != nil || (u.Data.SubURL == "" && trimStr == "") || (subURL.Scheme == "" && trimStr != "") {
		_, _ = u.SendMessage("Invalid URL. Please inspect your subURL or use [/url subURL] command once.", false)
		return
	} else {
		if trimStr != "" {
			u.Data.SubURL = subURL.String()
		}
		proxies, _, err := u.generateProxies(config.BotCfg.ConverterAPI)
		if err != nil {
			_ = u.EditMessage(u.MessageID, err.Error())
			return
		}

		var proxiesList []C.Proxy
		for _, v := range proxies {
			proxiesList = append(proxiesList, v)
		}
		// animation status
		go u.statusMessage("Retrieving IP information", checkFlag)

		start := time.Now()
		inbound, outbound := utils.GetIPList(proxiesList, config.BotCfg.MaxConn)
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
		_, _ = u.Bot.Send(warpFile)
		_ = u.DeleteMessage(u.MessageID)
	}
}

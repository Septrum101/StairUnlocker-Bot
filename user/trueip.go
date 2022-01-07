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

func (u *User) TrueIP(messageText string) {
	u.IsCheck = true
	checkFlag := make(chan bool)
	defer func() {
		u.IsCheck = false
		close(checkFlag)
	}()
	var subURL *url.URL
	subURL, err := url.Parse(strings.TrimSpace(strings.ReplaceAll(messageText, "/ip", "")))
	if err != nil || subURL.Scheme == "" {
		_, _ = u.Send("Invalid URL. Please inspect your subURL.", false)
		return
	} else {
		u.Data.SubURL = subURL.String()
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
		go func(u *User, checkFlag chan bool) {
			log.Infoln("[ID: %d] Retrieving IP information.", u.ID)
			count := 0
			for {
				select {
				case <-checkFlag:
					return
				default:
					count++
					if count > 5 {
						count = 0
					}
					_ = u.EditMessage(u.MessageID, fmt.Sprintf("Retrieving IP information%s", strings.Repeat(".", count)))
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(u, checkFlag)

		start := time.Now()
		inbound, outbound := utils.GetIPList(proxiesList, config.BotCfg.MaxConn)
		log.Warnln("[ID: %d] Total inbounds: %d -> outbounds: %d", u.ID, len(outbound), len(inbound))
		ipStatTitle := fmt.Sprintf("StairUnlocker Bot Bulletin:\nTotal %d nodes tested\nElapsed time: %s\ninbound IP: %d\noutbound IP: %d\nTimestamp: %s", len(proxies), time.Since(start).Round(time.Millisecond), len(outbound), len(inbound), time.Now().UTC().Format(time.RFC3339))
		ipStat := "StairUnlocker Bot Bulletin:\nEntrypoint IP: "
		for _, v := range inbound {
			ipStat += "\n" + v
		}
		ipStat += "\n\nEndpoint IP: "
		for _, v := range outbound {
			ipStat += "\n" + v
		}
		warpFile := tgBot.NewDocument(u.ID, tgBot.FileBytes{
			Name:  fmt.Sprintf("stairunlocker_bot_trueIP_%d.txt", time.Now().Unix()),
			Bytes: []byte(ipStat),
		})
		warpFile.Caption = fmt.Sprintf("%s\n%s\n@stairunlock_test_bot\nProject: https://git.io/Jyl5l", ipStatTitle, strings.Repeat("-", 25))
		_, _ = u.Bot.Send(warpFile)
		_ = u.DeleteMessage(u.MessageID)
	}
}

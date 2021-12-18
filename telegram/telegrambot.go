package telegram

import (
	"fmt"
	"github.com/Dreamacro/clash/log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/user"
	"net/url"
	"strings"
	"time"
)

func TGUpdates(buf *chan *user.User, userMap *map[int64]*user.User, cfg *config.SuConfig) (err error) {
	bot, err := tg.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		panic(err)
	}
	if cfg.LogLevel == 0 {
		bot.Debug = true
	}
	log.Infoln("Authorized on account %s", bot.Self.UserName)

	updateCfg := tg.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := bot.GetUpdatesChan(updateCfg)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// initial User struct
		usr, ok := (*userMap)[update.Message.Chat.ID]
		if !ok {
			usr = &user.User{
				ID:  update.Message.Chat.ID,
				Bot: bot,
			}
		}
		// select telegram cmd.
		switch {
		case update.Message.Text == "/start":
			_ = usr.Send(`
/url subURL - Test SubURL, Support http/https/vmess/ss/ssr/trojan.
/stat - Show the last checking result.
`)

		case update.Message.Text == "/stat":
			if usr.Data.CheckInfo == "" {
				_ = usr.Send("Cannot find the status information. Please use /url command first.")
			} else {
				_ = usr.Send((*userMap)[usr.ID].Data.CheckInfo)
			}

		case strings.HasPrefix(update.Message.Text, "/url"):
			var subURL *url.URL
			if len(*buf) > cfg.MaxOnline {
				_ = usr.Send("Too many connections, Please try again later.")
			} else {
				subURL, err = url.Parse(strings.TrimSpace(strings.ReplaceAll(update.Message.Text, "/url", "")))
				if err != nil || subURL.Scheme == "" {
					_ = usr.Send("Invalid URL")
				} else {
					// the time between previous testing.
					internal := time.Duration(cfg.Internal)
					if time.Since(time.Unix(usr.Data.LastCheck, 0)) < internal*time.Minute {
						remainTime := internal*time.Minute - time.Since(time.Unix(usr.Data.LastCheck, 0))
						_ = usr.Send(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)))
					} else {
						usr.Data = user.Data{LastCheck: time.Now().Unix(), SubURL: subURL.String()}
						*buf <- usr
						(*userMap)[update.Message.Chat.ID] = usr
						_ = usr.Send("Checking nodes status...")
					}
				}
			}

		default:
			_ = usr.Send("Invalid command")
		}
		log.Debugln("Telegram Bot: [ID: %d], Text: %s", update.Message.From.ID, update.Message.Text)
	}
	return
}
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

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		usr := user.User{
			ID:  update.Message.Chat.ID,
			Bot: bot,
		}

		if update.Message == nil {
			continue
		}
		// 选择telegram命令
		switch {
		case update.Message.Text == "/start":
			_ = usr.Send("/url Test subURL, Support vmess/ss/ssr/http/https.\n/stat Show last status.")

		case update.Message.Text == "/stat":
			if (*userMap)[usr.ID].Data.CheckInfo == "" {
				_ = usr.Send("Cannot find the status information. Please use /url command first.")
			} else {
				_ = usr.Send((*userMap)[usr.ID].Data.CheckInfo)
			}

		case strings.HasPrefix(update.Message.Text, "/url"):
			if len(*buf) > cfg.MaxOnline {
				_ = usr.Send(fmt.Sprintf("Too many connections, Please try again later."))
			} else {
				subURL, _ := url.Parse(strings.TrimSpace(strings.ReplaceAll(update.Message.Text, "/url", "")))
				if subURL.Scheme == "" {
					_ = usr.Send("Invalid URL")
				} else {
					schemeList := []string{"ss", "ssr", "trojan", "vemss"}
					for i := range schemeList {
						if subURL.Scheme == schemeList[i] {
							subURL.Fragment = strings.ReplaceAll(subURL.Fragment, "\n", "|")
						}
					}
					if exUser, ok := (*userMap)[update.Message.Chat.ID]; ok {
						// 用户测试间隔
						internal := time.Duration(cfg.Internal)
						if time.Now().Sub(time.Unix(exUser.Data.LastCheck, 0)) < internal*time.Minute {
							remainTime := internal*time.Minute - time.Now().Sub(time.Unix(exUser.Data.LastCheck, 0))
							_ = usr.Send(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)))
						} else {
							exUser.Data = user.Data{LastCheck: time.Now().Unix(), SubURL: subURL.String()}
							(*userMap)[update.Message.Chat.ID] = exUser
							_ = usr.Send("Checking nodes status...")
							*buf <- exUser
						}
					} else {
						usr.Data = user.Data{LastCheck: time.Now().Unix(), SubURL: subURL.String()}
						(*userMap)[update.Message.Chat.ID] = &usr
						_ = usr.Send("Checking nodes status...")
						*buf <- &usr
					}
				}
			}
		default:
			_ = usr.Send("Invalid command")
		}
		log.Debugln("TGUpdates Bot: [ID: %d], Text: %s", update.Message.From.ID, update.Message.Text)
	}
	return
}

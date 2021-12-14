package utils

import (
	"fmt"
	"github.com/Dreamacro/clash/log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"net/url"
	"strings"
	"time"
)

type TgBot struct {
	SendMessage string
	Check       string
	Bot         *tg.BotAPI
}

type User struct {
	ID   int64
	Data UserData
}

type UserData struct {
	LastCheck int64
	SubURL    string
	CheckInfo string
}

var conf *config.SuConfig

func (tb *TgBot) NewBot(cfg *config.SuConfig) {
	conf = cfg
	bot, err := tg.NewBotAPI(cfg.Telegram.TelegramToken)
	if err != nil {
		panic(err)
	}
	if cfg.LogLevel == 0 {
		bot.Debug = false
	}
	tb.Bot = bot
	log.Infoln("Authorized on account %s", bot.Self.UserName)
}

func (tb *TgBot) TelegramUpdates(buf *chan *User, userMap *map[int64]UserData) {
	bot := tb.Bot
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// 选择telegram命令
		switch {
		case update.Message.Text == "/start":
			_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "/check Check all node.\n/stat Show last status."))

		case update.Message.Text == "/stat":
			if (*userMap)[update.Message.Chat.ID].CheckInfo == "" {
				_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Cannot find the status information. Please use /url subURL ."))
			}
			_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, (*userMap)[update.Message.Chat.ID].CheckInfo))

		case strings.HasPrefix(update.Message.Text, "/url"):
			if len(*buf) > conf.MaxOnline {
				_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Too many connections, Please try again later.")))
			} else {
				subURL, _ := url.Parse(strings.TrimSpace(strings.ReplaceAll(update.Message.Text, "/url", "")))
				if subURL.Scheme == "" {
					_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Invalid URL"))
				} else {
					if usrData, ok := (*userMap)[update.Message.Chat.ID]; ok {
						// 每个用户10分钟只能测试一次
						if time.Now().Sub(time.Unix(usrData.LastCheck, 0)) < 10*time.Minute {
							remainTime := 10*time.Minute - time.Now().Sub(time.Unix(usrData.LastCheck, 0))
							_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Please retry after %s.", remainTime.Round(time.Second))))
						} else {
							_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Checking nodes status..."))
							*buf <- &User{
								ID:   update.Message.Chat.ID,
								Data: (*userMap)[update.Message.Chat.ID],
							}
						}
					} else {
						(*userMap)[update.Message.Chat.ID] = UserData{LastCheck: time.Now().Unix(), SubURL: subURL.String()}
						_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Checking nodes status..."))
						*buf <- &User{
							ID:   update.Message.Chat.ID,
							Data: (*userMap)[update.Message.Chat.ID],
						}
					}
				}
			}

		default:
			_, _ = bot.Send(tg.NewMessage(update.Message.Chat.ID, "Invalid command"))
		}
		log.Debugln("TelegramUpdates Bot: [ID: %d], Text: %s", update.Message.From.ID, update.Message.Text)
	}
}

package telegram

import (
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/user"
	"net/url"
	"strings"
	"time"
)

func Updates(buf *chan *user.User, userMap *map[int64]*user.User) (err error) {
	bot, err := tgBot.NewBotAPI(config.BotCfg.TelegramToken)
	if err != nil {
		panic(err)
	}
	if config.BotCfg.LogLevel == 0 {
		bot.Debug = true
	}
	log.Infoln("Authorized on account %s", bot.Self.UserName)

	updateCfg := tgBot.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := bot.GetUpdatesChan(updateCfg)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// If user is not exist, initial User struct
		usr, exist := (*userMap)[update.Message.Chat.ID]
		if !exist {
			(*userMap)[update.Message.Chat.ID] = &user.User{
				ID:  update.Message.Chat.ID,
				Bot: bot,
			}
			usr = (*userMap)[update.Message.Chat.ID]
		}
		// select telegram cmd.
		switch {
		case update.Message.Text == "/start":
			_ = usr.Send(`
/url subURL - Test SubURL, Support http/https/vmess/ss/ssr/trojan.
/retest - Retest last subURL.
/stat - Show the last checking result.
`)

		case update.Message.Text == "/stat":
			if usr.Data.CheckInfo == "" {
				_ = usr.Send("Cannot find status information. Please use /url subURL command first.")
			} else {
				_ = usr.Send((*userMap)[usr.ID].Data.CheckInfo)
			}

		case strings.HasPrefix(update.Message.Text, "/url") || update.Message.Text == "/retest":
			if len(*buf) > config.BotCfg.MaxOnline {
				_ = usr.Send("Too many connections, Please try again later.")
				continue
			}
			// limited double-checking
			if usr.IsCheck {
				_ = usr.Send("Duplication, Previous testing is not completed! Please try again later.")
				continue
			}
			if update.Message.Text == "/retest" {
				if usr.Data.SubURL == "" {
					_ = usr.Send("Cannot find subURL. Please use /url subURL command first.")
				} else if usr.UserOutInternal(config.BotCfg.Internal) {
					*buf <- usr
					(*userMap)[update.Message.Chat.ID].Data.LastCheck = time.Now().Unix()
				}
			} else {
				var subURL *url.URL
				replaceStrings := strings.NewReplacer("/url", "", "\r", "", "\n", "")
				rawStr := strings.TrimSpace(replaceStrings.Replace(update.Message.Text))
				subURL, err = url.Parse(rawStr)
				if err != nil || subURL.Scheme == "" {
					_ = usr.Send("Invalid URL. Please inspect your subURL.")
				} else if usr.UserOutInternal(config.BotCfg.Internal) {
					*buf <- usr
					(*userMap)[update.Message.Chat.ID].Data = user.Data{LastCheck: time.Now().Unix(), SubURL: subURL.String()}
				}
			}

		default:
			_ = usr.Send("Invalid command")
		}
		log.Debugln("Telegram Bot: [ID: %d], Text: %s", update.Message.From.ID, update.Message.Text)
	}
	return
}

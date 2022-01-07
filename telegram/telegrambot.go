package telegram

import (
	"fmt"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/user"
	"net/url"
	"strings"
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
	x, _ := bot.GetMyCommands()
	fmt.Println(x)
	_, _ = bot.Request(tgBot.SetMyCommandsConfig{Commands: []tgBot.BotCommand{
		{"url", "Test SubURL, Support http/https/vmess/ss/ssr/trojan."},
		{"ip", "Test Real IP information, Support http/https/vmess/ss/ssr/trojan."},
		{"stat", "Show the last checking result"},
	}})
	x, _ = bot.GetMyCommands()
	fmt.Println(x)
	// todo initial command list
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
			_, _ = usr.SendMessage(`
/url subURL - Test SubURL, Support http/https/vmess/ss/ssr/trojan.
/ip subURL - Test Real IP information, Support http/https/vmess/ss/ssr/trojan.
/stat - Show the last checking result.
`, false)

		case update.Message.Text == "/stat":
			if usr.Data.CheckInfo == "" {
				_, _ = usr.SendMessage("Cannot find status information. Please use /url subURL command once.", false)
			} else {
				_, _ = usr.SendMessage((*userMap)[usr.ID].Data.CheckInfo, false)
			}

		case strings.HasPrefix(update.Message.Text, "/url") || strings.HasPrefix(update.Message.Text, "/ip"):
			// delete user privacy info
			_ = usr.DeleteMessage(update.Message.MessageID)
			if len(*buf) > config.BotCfg.MaxOnline {
				_, _ = usr.SendMessage("Too many connections, Please try again later.", false)
				continue
			}
			// forbid double-checking
			if usr.IsCheck {
				_, _ = usr.SendMessage("Duplication, Previous testing is not completed! Please try again later.", false)
				continue
			}

			var subURL *url.URL
			if strings.HasPrefix(update.Message.Text, "/url") {
				trimStr := strings.TrimSpace(strings.ReplaceAll(update.Message.Text, "/url", ""))
				subURL, err = url.Parse(trimStr)
				if err != nil || (usr.Data.SubURL == "" && trimStr == "") {
					_, _ = usr.SendMessage("Invalid URL. Please inspect your subURL or use /url subURL command once.", false)
				} else if usr.UserOutInternal(config.BotCfg.Internal) {
					if trimStr != "" {
						usr.Data.SubURL = subURL.String()
					}
					*buf <- usr
				}
			} else {
				if usr.UserOutInternal(config.BotCfg.Internal) {
					go usr.RealIP(update.Message.Text)
				}
			}

		default:
			_, _ = usr.SendMessage("Invalid command", false)
		}
		log.Debugln("Telegram Bot: [ID: %d], Text: %s", update.Message.From.ID, update.Message.Text)
	}
	return
}

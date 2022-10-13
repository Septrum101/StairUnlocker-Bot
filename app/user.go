package app

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/thank243/StairUnlocker-Bot/model"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

type User struct {
	ID              int64
	message         chan *tg.Message
	s               *Server
	editMsgID       int
	IsCheck         bool
	RefuseMessageID int
	Data            struct {
		LastCheck int64
		SubURL    string
		CheckInfo string
	}
}

func NewUser(server *Server, up *tg.Update) *User {
	u := &User{
		ID:      up.Message.Chat.ID,
		s:       server,
		message: make(chan *tg.Message),
	}

	go u.listenMessage()

	return u
}

func (u *User) listenMessage() {
	for msg := range u.message {
		// delete user privacy info
		u.DeleteMessage(msg.MessageID)

		switch {
		case msg.Text == "/start":
			u.cmdStart()
		case msg.Text == "/stat":
			u.cmdStat()
		case msg.Text == "/version":
			u.cmdVersion()
		case strings.HasPrefix(msg.Text, "/url"):
			if !u.validator() {
				continue
			}
			if err := u.cmdURL(msg.Text); err != nil {
				continue
			}
		case strings.HasPrefix(msg.Text, "/ip"):
			if !u.validator() {
				continue
			}
			if err := u.cmdIP(msg.Text); err != nil {
				continue
			}
		default:
			u.SendMessage("Invalid command")
		}
		log.Debugln("Telegram Bot: [ID: %d], Text: %s", msg.From.ID, msg.Text)
	}
}

func (u *User) cmdStart() {
	cmdList, _ := u.s.Bot.GetMyCommands()
	var str string
	for i := range cmdList {
		str += fmt.Sprintf("/%s - %s\n", cmdList[i].Command, cmdList[i].Description)
	}
	str += "The bot will use latest subURL for testing after a valid subURL."
	u.SendMessage(str)
}

func (u *User) cmdStat() {
	if u.Data.CheckInfo == "" {
		u.SendMessage("Cannot find status information. Please use [/url subURL] command once.")
	} else {
		u.SendMessage(u.s.userMap[u.ID].Data.CheckInfo)
	}
}

func (u *User) cmdVersion() {
	todayUser := 0
	for _, v := range u.s.userMap {
		if time.Now().Unix()-v.Data.LastCheck < int64(24*time.Hour.Seconds()) {
			todayUser++
		}
	}
	uptime := utils.FormatTime(time.Since(u.s.StartTime))
	u.SendMessage(fmt.Sprintf("StairUnlocker Bot %s\nUsers: (%d/%d) \nUptime: %s", C.Version, todayUser, len(u.s.userMap), uptime))

}

func (u *User) validator() bool {
	if len(u.s.userMap) > model.BotCfg.MaxOnline {
		u.SendMessage("Too many connections, Please try again later.")
		return false
	}
	// forbid double-checking
	if u.IsCheck {
		u.SendMessage("Duplication, Previous testing is not completed! Please try again later.")
		return false
	}

	return true
}

func (u *User) cmdURL(msg string) error {
	subURL, err := url.Parse(strings.TrimSpace(strings.ReplaceAll(msg, "/url", "")))
	if err != nil || (u.Data.SubURL == "" && subURL.String() == "") {
		u.SendMessage("Invalid URL. Please inspect your subURL or use [/url subURL] command once.")
		return err
	}
	if u.UserOutInternal() {
		u.streamMedia(subURL.String())
	}

	return nil
}

func (u *User) cmdIP(msg string) error {
	subURL, err := url.Parse(strings.TrimSpace(strings.ReplaceAll(msg, "/ip", "")))
	if err != nil || (u.Data.SubURL == "" && subURL.String() == "") {
		u.SendMessage("Invalid URL. Please inspect your subURL or use [/ip subURL] command once.")
		return err
	}
	if u.UserOutInternal() {
		u.realIP(subURL.String())
	}

	return nil
}

func (u *User) SendMessage(msg string) (resp tg.Message, err error) {
	resp, err = u.s.Bot.Send(tg.NewMessage(u.ID, msg))
	if err != nil {
		return
	}
	u.editMsgID = resp.MessageID
	return
}

func (u *User) UserOutInternal() bool {
	internal := time.Duration(model.BotCfg.Internal)
	if remainTime := internal*time.Second - time.Since(time.Unix(u.Data.LastCheck, 0)); remainTime > 0 {
		if u.RefuseMessageID == 0 {
			resp, _ := u.SendMessage(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)))
			u.RefuseMessageID = resp.MessageID
			go func() {
				n := 5 * time.Second
				for {
					remainTime := internal*time.Second - time.Since(time.Unix(u.Data.LastCheck, 0))
					if remainTime <= 0*time.Second {
						_ = u.DeleteMessage(u.RefuseMessageID)
						u.RefuseMessageID = 0
						return
					} else {
						_ = u.EditMessage(u.RefuseMessageID, fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)))
					}
					if remainTime <= 10*time.Second {
						n = 500 * time.Millisecond
					}
					time.Sleep(n)
				}
			}()
		}
		return false
	} else {
		return true
	}
}

func (u *User) DeleteMessage(msgID int) error {
	_, err := u.s.Bot.Send(tg.NewDeleteMessage(u.ID, msgID))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EditMessage(msgID int, text string) error {
	_, err := u.s.Bot.Send(tg.NewEditMessageText(u.ID, msgID, text))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) loading(info string, checkFlag chan bool) {
	log.Debugln("[ID: %d] %s", u.ID, info)
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
			u.EditMessage(u.editMsgID, fmt.Sprintf("%s%s", info, strings.Repeat(".", count)))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

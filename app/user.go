package app

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dreamacro/clash/log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/atomic"

	"github.com/thank243/StairUnlocker-Bot/model"
)

type User struct {
	ID             int64
	message        chan *tg.Message
	s              *Server
	isCheck        atomic.Bool
	countDownMsgID atomic.Int64
	data           struct {
		lastCheck   atomic.Int64
		subURL      atomic.String
		checkedInfo atomic.String
	}
	l sync.RWMutex
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
			if u.validator() {
				go u.cmdURL(msg.Text)
			}
		case strings.HasPrefix(msg.Text, "/ip"):
			if u.validator() {
				go u.cmdIP(msg.Text)
			}
		default:
			u.SendMessage("Invalid command")
		}
		log.Debugln("Telegram Bot: [ID: %d], Text: %s", msg.From.ID, msg.Text)
	}
}

func (u *User) validator() bool {
	if u.s.runningTask.Load() > int64(model.BotCfg.MaxOnline) {
		u.SendMessage("Too many connections, Please try again later.")
		return false
	}
	// forbid double-checking
	if u.isCheck.Load() {
		u.SendMessage("Duplication, Previous testing is not completed! Please try again later.")
		return false
	}

	return true
}

func (u *User) SendMessage(msg string) (resp tg.Message, err error) {
	resp, err = u.s.Bot.Send(tg.NewMessage(u.ID, msg))
	if err != nil {
		return
	}
	return
}

func (u *User) rateLimiting() bool {
	remainTime := float64(model.BotCfg.Internal) - time.Since(time.Unix(u.data.lastCheck.Load(), 0)).Seconds()
	if remainTime > 0 {
		if u.countDownMsgID.Load() == 0 {
			resp, _ := u.SendMessage(fmt.Sprintf("Please try again after %ds.", int(remainTime)))
			u.countDownMsgID.Store(int64(resp.MessageID))
			go func() {
				n := 5 * time.Second
				for {
					remainTime = float64(model.BotCfg.Internal) - time.Since(time.Unix(u.data.lastCheck.Load(), 0)).Seconds()
					if remainTime <= 0 {
						_ = u.DeleteMessage(int(u.countDownMsgID.Load()))
						u.countDownMsgID.Store(0)
						return
					} else {
						_ = u.EditMessage(int(u.countDownMsgID.Load()), fmt.Sprintf("Please try again after %ds.", int(remainTime)))
					}
					if remainTime <= 10 {
						n = 500 * time.Millisecond
					}
					time.Sleep(n)
				}
			}()
		}
		return true
	} else {
		return false
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
	u.l.Lock()
	_, err := u.s.Bot.Send(tg.NewEditMessageText(u.ID, msgID, text))
	u.l.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) loading(info string, checkFlag chan bool, msgID int) {
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
			u.EditMessage(msgID, fmt.Sprintf("%s%s", info, strings.Repeat(".", count)))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

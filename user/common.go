package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dreamacro/clash/log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *User) SendMessage(ctx string, isKeepSession bool) (resp tg.Message, err error) {
	resp, err = u.Bot.Send(tg.NewMessage(u.ID, ctx))
	if err != nil {
		return
	}
	if isKeepSession {
		u.MessageID = resp.MessageID
	}
	return
}

func (u *User) UserOutInternal(n int) bool {
	internal := time.Duration(n)
	if remainTime := internal*time.Second - time.Since(time.Unix(u.Data.LastCheck, 0)); remainTime > 0 {
		if u.RefuseMessageID == 0 {
			resp, _ := u.SendMessage(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)), true)
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
	_, err := u.Bot.Send(tg.NewDeleteMessage(u.ID, msgID))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EditMessage(msgID int, text string) error {
	_, err := u.Bot.Send(tg.NewEditMessageText(u.ID, msgID, text))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) statusMessage(info string, checkFlag chan bool) {
	log.Infoln("[ID: %d] %s", u.ID, info)
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
			_ = u.EditMessage(u.MessageID, fmt.Sprintf("%s%s", info, strings.Repeat(".", count)))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

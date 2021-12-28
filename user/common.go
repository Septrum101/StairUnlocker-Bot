package user

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (u *User) Send(ctx string, isSaveMessageID bool) (resp tg.Message, err error) {
	resp, err = u.Bot.Send(tg.NewMessage(u.ID, ctx))
	if err != nil {
		return
	}
	if isSaveMessageID {
		u.MessageID = resp.MessageID
	}
	return
}

func (u *User) UserOutInternal(n int) bool {
	internal := time.Duration(n)
	if remainTime := internal*time.Minute - time.Since(time.Unix(u.Data.LastCheck, 0)); remainTime > 0 {
		if u.RefuseMessageID == 0 {
			resp, _ := u.Send(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)), true)
			u.RefuseMessageID = resp.MessageID
			go func() {
				n := 5 * time.Second
				for {
					remainTime := internal*time.Minute - time.Since(time.Unix(u.Data.LastCheck, 0))
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

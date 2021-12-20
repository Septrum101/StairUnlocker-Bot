package user

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (u *User) Send(ctx string) (resp tg.Message, err error) {
	resp, err = u.Bot.Send(tg.NewMessage(u.ID, ctx))
	if err != nil {
		return
	}
	u.MessageID = resp.MessageID
	return
}

func (u *User) UserOutInternal(n int) bool {
	internal := time.Duration(n)
	if time.Since(time.Unix(u.Data.LastCheck, 0)) < internal*time.Minute {
		remainTime := internal*time.Minute - time.Since(time.Unix(u.Data.LastCheck, 0))
		_, _ = u.Send(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)))
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

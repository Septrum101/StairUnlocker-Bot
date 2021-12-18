package user

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func (u *User) Send(ctx string) (err error) {
	_, err = u.Bot.Send(tg.NewMessage(u.ID, ctx))
	if err != nil {
		return
	}
	return
}

func (u *User) UserOutInternal(n int) bool {
	internal := time.Duration(n)
	if time.Since(time.Unix(u.Data.LastCheck, 0)) < internal*time.Minute {
		remainTime := internal*time.Minute - time.Since(time.Unix(u.Data.LastCheck, 0))
		_ = u.Send(fmt.Sprintf("Please try again after %s.", remainTime.Round(time.Second)))
		return false
	} else {
		return true
	}
}

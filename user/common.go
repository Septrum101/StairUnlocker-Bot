package user

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *User) Send(ctx string) (err error) {
	_, err = u.Bot.Send(tg.NewMessage(u.ID, ctx))
	if err != nil {
		return
	}
	return
}

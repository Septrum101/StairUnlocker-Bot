package user

import tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type User struct {
	ID              int64
	Data            Data
	Bot             *tgBot.BotAPI
	IsCheck         bool
	MessageID       int
	RefuseMessageID int
}

type Data struct {
	LastCheck int64
	SubURL    string
	CheckInfo string
}

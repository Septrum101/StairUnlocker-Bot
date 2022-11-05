package app

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type Server struct {
	userMap        map[int64]*User
	updatesChannel tgBot.UpdatesChannel
	StartTime      time.Time
	Bot            *tgBot.BotAPI
	l              sync.RWMutex
	runningTask    atomic.Int64
}

func NewServer() (*Server, error) {
	bot, err := tgBot.NewBotAPI(viper.GetString("telegramToken"))
	if err != nil {
		return nil, err
	}
	if viper.GetString("log_level") == "debug" {
		bot.Debug = true
	}
	log.Infoln("Authorized on account %s", bot.Self.UserName)
	// initial command list
	preCommands, _ := bot.GetMyCommands()
	currCommands := []tgBot.BotCommand{
		{"url", "Get nodes unlock status."},
		{"ip", "Get Real IP information."},
		{"stat", "Show the latest checking result."},
		{"version", "Show version."},
	}
	if fmt.Sprint(preCommands) != fmt.Sprint(currCommands) {
		_, err = bot.Request(tgBot.SetMyCommandsConfig{Commands: currCommands})
		if err != nil {
			log.Errorln(err.Error())
		}
	}

	updateCfg := tgBot.NewUpdate(0)
	updateCfg.Timeout = 60
	s := &Server{
		updatesChannel: bot.GetUpdatesChan(updateCfg),
		StartTime:      time.Now(),
		userMap:        make(map[int64]*User),
		Bot:            bot,
	}

	return s, nil
}

func (s *Server) Start() {
	for i := range s.updatesChannel {
		if i.Message == nil || i.Message.Text == "" {
			continue
		}
		if u, ok := s.userMap[i.SentFrom().ID]; ok {
			u.message <- i.Message
		} else {
			u = NewUser(s, &i)
			s.userMap[i.SentFrom().ID] = u
			u.message <- i.Message
		}
	}
}

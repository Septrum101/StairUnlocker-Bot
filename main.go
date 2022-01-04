package main

import (
	"flag"
	"fmt"
	"github.com/Dreamacro/clash/log"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/telegram"
	"github.com/thank243/StairUnlocker-Bot/user"
)

func main() {
	//command-line
	if config.Version {
		return
	}
	if config.Help {
		flag.PrintDefaults()
		return
	}
	log.SetLevel(config.BotCfg.LogLevel)
	fmt.Printf("Log Level: %s\n", config.BotCfg.LogLevel)

	// receive the user from telegram
	ch := make(chan *user.User, config.BotCfg.MaxOnline)
	userMap := make(map[int64]*user.User)
	go func() {
		err := telegram.Updates(&ch, &userMap)
		if err != nil {
			log.Errorln(err.Error())
		}
	}()
	for usr := range ch {
		go usr.URLCheck()
	}
}

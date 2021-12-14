package main

import (
	"flag"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
	"gopkg.in/yaml.v3"
	"runtime"
	"strings"
	"time"
)

var (
	suConfig    *config.SuConfig
	tg          utils.TgBot
	proxiesList []C.Proxy
	start       time.Time
	ver         bool
	help        bool
	configFile  string
	UserMap     map[int64]utils.UserData
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&ver, "v", false, "show current ver of StairUnlock")
	flag.StringVar(&configFile, "f", "", "specify configuration file")
	flag.Parse()
	suConfig = config.Init(&configFile)
}

func URLCheck(user utils.User) {
	proxies, unmarshalProxies, err := config.GenerateProxies(suConfig.ConverterAPI, UserMap[user.ID].SubURL)
	if err != nil {
		_, _ = tg.Bot.Send(tgBot.NewMessage(user.ID, err.Error()))
		return
	}
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}

	//同时连接数
	connNum := suConfig.MaxConn
	if i := len(proxiesList); i < connNum {
		connNum = i
	}
	start = time.Now()
	netflixList := utils.BatchCheck(proxiesList, connNum)
	report := fmt.Sprintf("Total %d nodes test completed, %d unlock nodes, Elapsed time: %s", len(proxiesList), len(netflixList), time.Now().Sub(start).Round(time.Millisecond))

	log.Warnln(report)
	telegramReport := fmt.Sprintf("%s, Timestamp: %s", report, time.Now().Round(time.Millisecond))
	UserMap[user.ID] = utils.UserData{
		LastCheck: UserMap[user.ID].LastCheck,
		SubURL:    UserMap[user.ID].SubURL,
		CheckInfo: telegramReport,
	}
	marshal, _ := yaml.Marshal(config.NETFLIXFilter(netflixList, unmarshalProxies))
	_, _ = tg.Bot.Send(tgBot.NewMessage(user.ID, fmt.Sprintf("%s\n%s\n%s", telegramReport, strings.Repeat("-", 50), string(marshal))))
}

func main() {
	versionStr := fmt.Sprintf("StairUnlock-Bot %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime)
	//command-line
	if ver {
		fmt.Printf(versionStr)
		return
	}
	if help {
		fmt.Printf(versionStr)
		flag.PrintDefaults()
		return
	}
	fmt.Printf(versionStr)
	log.SetLevel(suConfig.LogLevel)
	fmt.Printf("Log Level: %s\n", suConfig.LogLevel)

	// 初始化telegramBot
	tg.NewBot(suConfig)
	ch := make(chan *utils.User, suConfig.MaxOnline)
	UserMap = make(map[int64]utils.UserData)
	go func() { tg.TelegramUpdates(&ch, &UserMap) }()

	for user := range ch {
		go func(usr utils.User) { URLCheck(usr) }(*user)
	}
}

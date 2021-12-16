package main

import (
	"flag"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/telegram"
	"github.com/thank243/StairUnlocker-Bot/user"
	"runtime"
)

var (
	SuConfig   *config.SuConfig
	ver        bool
	help       bool
	configPath string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&ver, "v", false, "show current ver of StairUnlock")
	flag.StringVar(&configPath, "f", "", "specify configuration file")
	flag.Parse()
	SuConfig = config.Init(&configPath)
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
	log.SetLevel(SuConfig.LogLevel)
	fmt.Printf("Log Level: %s\n", SuConfig.LogLevel)

	ch := make(chan *user.User, SuConfig.MaxOnline)
	userMap := make(map[int64]user.User)
	go func() { _ = telegram.TGUpdates(&ch, &userMap, SuConfig) }()

	for u := range ch {
		go func(u *user.User) { u.URLCheck(SuConfig.ConverterAPI, SuConfig.MaxConn) }(u)
	}
}

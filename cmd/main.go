package main

import (
	"flag"
	"fmt"
	"runtime"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"

	"github.com/thank243/StairUnlocker-Bot/app"
	"github.com/thank243/StairUnlocker-Bot/model"
)

func main() {
	// command-line
	if model.Version {
		fmt.Printf("StairUnlock-Bot %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime)
		return
	}

	if model.Help {
		flag.PrintDefaults()
		return
	}

	log.SetLevel(model.BotCfg.LogLevel)
	fmt.Printf("Log Level: %s\n", model.BotCfg.LogLevel)

	s, err := app.NewServer()
	if err != nil {
		log.Fatalln("%v", err)
	}
	s.Start()
}

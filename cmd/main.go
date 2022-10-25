package main

import (
	"flag"
	"fmt"
	"runtime"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/spf13/viper"

	"github.com/thank243/StairUnlocker-Bot/app"
)

func main() {
	fmt.Printf(fmt.Sprintf("StairUnlock-Bot %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime))

	// command-line
	var (
		Version bool
		Help    bool
	)

	flag.BoolVar(&Help, "h", false, "this help")
	flag.BoolVar(&Version, "v", false, "show current version of StairUnlock")
	flag.Parse()

	if Version {
		fmt.Printf("StairUnlock-Bot %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime)
		return
	}

	if Help {
		flag.PrintDefaults()
		return
	}

	log.SetLevel(log.LogLevelMapping[viper.GetString("log_level")])
	fmt.Printf("Log Level: %s\n", viper.GetString("log_level"))

	s, err := app.NewServer()
	if err != nil {
		log.Fatalln("%v", err)
	}
	s.Start()
}

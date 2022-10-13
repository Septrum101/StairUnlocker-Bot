package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/model"
)

func init() {
	fmt.Printf(fmt.Sprintf("StairUnlock-Bot %s %s %s with %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), C.BuildTime))

	flag.BoolVar(&model.Help, "h", false, "this help")
	flag.BoolVar(&model.Version, "v", false, "show current version of StairUnlock")
	flag.StringVar(&model.ConfPath, "f", "", "specify configuration file")
	flag.Parse()

	//initial config.yaml
	var (
		buf []byte
	)

	if model.ConfPath != "" {
		buf, _ = ioutil.ReadFile(model.ConfPath)
	} else {
		_, err := os.Stat("config.yaml")
		if err != nil {
			b, _ := ioutil.ReadFile("config.example.yaml")
			ioutil.WriteFile("config.yaml", b, 644)
		}
		buf, _ = ioutil.ReadFile("config.yaml")
	}
	yaml.Unmarshal(buf, &model.BotCfg)
}

func UnmarshalRawConfig(buf []byte) (*model.RawConfig, error) {
	rawConf := &model.RawConfig{}
	if err := yaml.Unmarshal(buf, rawConf); err != nil {
		return nil, err
	}
	return rawConf, nil
}

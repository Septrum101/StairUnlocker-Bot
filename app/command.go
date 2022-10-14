package app

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"

	"github.com/thank243/StairUnlocker-Bot/utils"
)

func (u *User) cmdStart() {
	cmdList, _ := u.s.Bot.GetMyCommands()
	var str string
	for i := range cmdList {
		str += fmt.Sprintf("/%s - %s\n", cmdList[i].Command, cmdList[i].Description)
	}
	str += "The bot will use latest subURL for testing after a valid subURL."
	u.SendMessage(str)
}

func (u *User) cmdStat() {
	if u.data.checkedInfo.Load() == "" {
		u.SendMessage("Cannot find status information. Please use [/url subURL] command once.")
	} else {
		u.SendMessage(u.s.userMap[u.ID].data.checkedInfo.Load())
	}
}

func (u *User) cmdVersion() {
	todayUser := 0
	for _, v := range u.s.userMap {
		if time.Now().Unix()-v.data.lastCheck.Load() < int64(24*time.Hour.Seconds()) {
			todayUser++
		}
	}
	uptime := utils.FormatTime(time.Since(u.s.StartTime))
	u.SendMessage(fmt.Sprintf("StairUnlocker Bot %s\nUsers: (%d/%d) \nUptime: %s", C.Version, todayUser, len(u.s.userMap), uptime))

}

func (u *User) cmdURL(msg string) error {
	subURL, err := url.Parse(strings.TrimSpace(strings.ReplaceAll(msg, "/url", "")))
	if err != nil || (u.data.subURL.Load() == "" && subURL.String() == "") {
		u.SendMessage("Invalid URL. Please inspect your subURL or use [/url subURL] command once.")
		return err
	}

	su := subURL.String()
	if su == "" {
		su = u.data.subURL.Load()
	}
	u.s.runningTask.Add(1)
	u.streamMedia(su)
	u.s.runningTask.Add(-1)

	return nil
}

func (u *User) cmdIP(msg string) error {
	subURL, err := url.Parse(strings.TrimSpace(strings.ReplaceAll(msg, "/ip", "")))
	if err != nil || (u.data.subURL.Load() == "" && subURL.String() == "") {
		u.SendMessage("Invalid URL. Please inspect your subURL or use [/ip subURL] command once.")
		return err
	}

	su := subURL.String()
	if su == "" {
		su = u.data.subURL.Load()
	}
	u.s.runningTask.Add(1)
	u.realIP(su)
	u.s.runningTask.Add(-1)

	return nil
}

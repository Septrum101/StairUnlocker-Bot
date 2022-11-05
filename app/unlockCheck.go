package app

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"

	"github.com/thank243/StairUnlocker-Bot/model"
	"github.com/thank243/StairUnlocker-Bot/provider"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

var l sync.RWMutex

func statistic(streamMediaList *[]model.StreamData) map[string]int {
	statMap := make(map[string]int)
	for i := range *streamMediaList {
		statMap[(*streamMediaList)[i].Name]++
		if !(*streamMediaList)[i].Unlock {
			statMap[(*streamMediaList)[i].Name]--
		}
	}
	return statMap
}

func (u *User) streamMedia(subUrl string) error {
	u.isCheck.Store(true)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		u.data.lastCheck.Store(time.Now().Unix())
		u.isCheck.Store(false)
		cancel()
	}()

	msgInst, _ := u.SendMessage("Converting from API server.")
	proxies, err := u.buildProxies(subUrl)
	if err != nil {
		u.EditMessage(msgInst.MessageID, err.Error())
		return err
	}

	// animation while waiting test.
	go u.loading(ctx, "Checking nodes unlock status", msgInst.MessageID)

	var proxiesList []C.Proxy
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	// Must have valid node.
	if len(proxiesList) > 0 {
		log.Infoln("[ID: %d] Start unlock test", u.ID)
		start := time.Now()
		unlockList := u.batch(proxiesList, viper.GetInt("maxConn"))
		cancel()

		report := fmt.Sprintf("Total %d nodes, Duration: %s", len(proxiesList), time.Since(start).Round(time.Millisecond))

		var nameList []string
		statisticMap := statistic(&unlockList)
		for k := range statisticMap {
			nameList = append(nameList, k)
		}
		sort.Strings(nameList)
		var finalStr string
		for i := range nameList {
			finalStr += fmt.Sprintf("%s: %d\n", nameList[i], statisticMap[nameList[i]])
		}
		telegramReport := fmt.Sprintf("StairUnlocker Bot %s Bulletin:\n%s\n%sTimestamp: %s\n%s", C.Version, report, finalStr, time.Now().UTC().Format(time.RFC3339), strings.Repeat("-", 25))
		log.Warnln("[ID: %d] %s", u.ID, report)
		u.EditMessage(msgInst.MessageID, "Uploading PNG file...")

		buffer, err := utils.GeneratePNG(unlockList, nameList)
		if err != nil {
			return err
		}
		// send result image
		wrapPNG := tgBot.NewDocument(u.ID, tgBot.FileBytes{
			Name:  fmt.Sprintf("stairunlocker_bot_result_%d.png", time.Now().Unix()),
			Bytes: buffer.Bytes(),
		})
		wrapPNG.Caption = fmt.Sprintf("%s\n@stairunlock_test_bot\nProject: https://git.io/Jyl5l", telegramReport)
		// save test results.
		u.data.checkedInfo.Store(wrapPNG.Caption)

		u.s.Bot.Send(wrapPNG)
		u.DeleteMessage(msgInst.MessageID)
	}
	return nil
}

// batch : n int, to set ConcurrencyNum.
func (u *User) batch(proxiesList []C.Proxy, n int) (streamDataList []model.StreamData) {
	type combineProxy struct {
		proxy  C.Proxy
		stream provider.AbsStream
	}

	var (
		wg sync.WaitGroup
		cp []combineProxy
	)

	streamList := provider.NewStreamList()
	for i := range proxiesList {
		for ii := range streamList {
			cp = append(cp, combineProxy{
				proxy:  proxiesList[i],
				stream: streamList[ii],
			})
		}
	}

	streamDataList = make([]model.StreamData, 0, len(streamList)*len(proxiesList))

	// prefix for node name on log.
	curr, total := int32(0), len(cp)
	// initial pool
	pool, err := ants.NewPoolWithFunc(n, func(i interface{}) {
		c := i.(combineProxy)
		result, err := c.stream.IsUnlock(&c.proxy)

		rtt, _ := time.ParseDuration(result.Latency)
		if rtt > 10*time.Second {
			log.Warnln("slow running: %s %v %s", result.ProxyName, rtt, u.data.subURL.Load())
		}

		atomic.AddInt32(&curr, 1)
		if err != nil {
			log.Debugln("(%d/%d) %s : %s", atomic.LoadInt32(&curr), total, c.proxy.Name(), err.Error())
		} else if result.Unlock {
			log.Debugln("(%d/%d) %s | %s Unlock", atomic.LoadInt32(&curr), total, c.proxy.Name(), result.Name)
		} else {
			log.Debugln("(%d/%d) %s | %s None", atomic.LoadInt32(&curr), total, c.proxy.Name(), result.Name)
		}
		l.Lock()
		streamDataList = append(streamDataList, result)
		l.Unlock()
		wg.Done()
	})
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	defer pool.Release()

	for i := range cp {
		wg.Add(1)
		err = pool.Invoke(cp[i])
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	wg.Wait()
	return
}

package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/go-resty/resty/v2"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/panjf2000/ants/v2"

	"github.com/thank243/StairUnlocker-Bot/model"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

type geoIP struct {
	Country string `json:"country"`
	Isp     string `json:"isp"`
	Query   string `json:"query"`
	Status  string `json:"status"`
}

const ipAPI = "http://ip-api.com/json?fields=25089"
const ipAPIBatch = "http://ip-api.com/batch?fields=25089"

func (u *User) realIP(subUrl string) error {
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
	if subUrl != "" {
		u.data.subURL.Store(subUrl)
	}

	// animation while waiting test.
	go u.loading(ctx, "Retrieving IP information", msgInst.MessageID)

	var proxiesList []C.Proxy
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}
	// Must have valid node.
	if len(proxiesList) > 0 {
		log.Infoln("[ID: %d] Start real IP test", u.ID)
		start := time.Now()
		inbound, outbound := getIPList(proxiesList, model.BotCfg.MaxConn)
		cancel()

		duration := time.Since(start).Round(time.Millisecond)
		log.Warnln("[ID: %d] Total %d nodes: inbounds: %d -> outbounds: %d, Duration: %s", u.ID, len(proxiesList), len(inbound), len(outbound), duration)
		ipStatTitle := fmt.Sprintf("StairUnlocker Bot %s Bulletin:\nTotal %d nodes, Duration: %s\ninbound IP: %d\noutbound IP: %d\nTimestamp: %s", C.Version, len(proxies), duration, len(inbound), len(outbound), time.Now().UTC().Format(time.RFC3339))
		ipStat := fmt.Sprintf("StairUnlocker Bot %s Bulletin:\nEntrypoint IP: ", C.Version)
		for _, v := range inbound {
			ipStat += "\n" + v
		}
		ipStat += "\n\nEndpoint IP: "
		for i := range outbound {
			ipStat += "\n" + outbound[i]
		}
		warpFile := tgBot.NewDocument(u.ID, tgBot.FileBytes{
			Name:  fmt.Sprintf("stairunlocker_bot_realIP_%d.txt", time.Now().Unix()),
			Bytes: []byte(ipStat),
		})
		warpFile.Caption = fmt.Sprintf("%s\n%s\n@stairunlock_test_bot\nProject: https://git.io/Jyl5l", ipStatTitle, strings.Repeat("-", 25))
		u.s.Bot.Send(warpFile)
		u.DeleteMessage(msgInst.MessageID)
	}
	return nil
}

func endIPTest(p C.Proxy) (ipInfo geoIP, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr, err := utils.UrlToMetadata(ipAPI)
	if err != nil {
		return
	}

	instance, err := p.DialContext(ctx, &addr)
	if err != nil {
		return
	}

	defer func(instance C.Conn) {
		err = instance.Close()
		if err != nil {
			return
		}
	}(instance)

	_, err = resty.New().
		SetTransport(&http.Transport{
			DialContext: func(context.Context, string, string) (net.Conn, error) {
				return instance, err
			}}).
		SetCloseConnection(true).SetHeader("User-Agent", "curl").
		R().SetContext(ctx).SetResult(&ipInfo).Get(ipAPI)
	return
}

func entryIPTest(proxiesList []C.Proxy) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ipInfos []string
	var ipList []string
	for i := range proxiesList {
		ipList = append(ipList, entryIP(proxiesList[i].Addr())...)
	}

	offset := 0
	for {
		var ips []geoIP
		end := offset + 100
		if offset+100 > len(proxiesList) {
			end = len(proxiesList)
		}
		_, err := resty.New().SetHeader("User-Agent", "curl").SetRetryCount(3).
			SetCloseConnection(true).R().SetContext(ctx).SetResult(&ips).SetBody(ipList[offset:end]).
			Post(ipAPIBatch)
		if err != nil {
			return nil, err
		}
		for i := range ips {
			if ips[i].Status == "success" {
				ipInfos = append(ipInfos, fmt.Sprintf("%s - %s, ISP: %s", ips[i].Query, ips[i].Country, ips[i].Isp))
			}
		}
		if end == len(proxiesList) {
			break
		}
		offset += 100
	}
	return ipInfos, nil
}

func entryIP(addr string) []string {
	addrList, _ := net.LookupHost(strings.Split(addr, ":")[0])
	return addrList
}

func deDuplication(list []string) []string {
	countMap := make(map[string]int)
	for i := range list {
		countMap[list[i]]++
	}
	var countList []string
	for k, v := range countMap {
		countList = append(countList, fmt.Sprintf("%s (%d)", k, v))
	}
	sort.Strings(countList)
	return countList
}

func getIPList(proxiesList []C.Proxy, n int) ([]string, []string) {
	var (
		wg          sync.WaitGroup
		endIPList   []string
		entryIPList []string
	)

	entryIPList, err := entryIPTest(proxiesList)
	if err != nil {
		log.Errorln("%v", err)
	}

	pool, err := ants.NewPoolWithFunc(n, func(i interface{}) {
		p := i.(C.Proxy)
		resp, _ := endIPTest(p)
		if resp.Query != "" {
			l.Lock()
			endIPList = append(endIPList, fmt.Sprintf("%s - %s, ISP: %s", resp.Query, resp.Country, resp.Isp))
			l.Unlock()
		}
		wg.Done()
	})
	defer pool.Release()
	if err != nil {
		log.Errorln(err.Error())
		return nil, nil
	}

	for i := range proxiesList {
		wg.Add(1)
		err = pool.Invoke(proxiesList[i])
		if err != nil {
			log.Errorln(err.Error())
			return nil, nil
		}
	}
	wg.Wait()
	return deDuplication(entryIPList), deDuplication(endIPList)
}

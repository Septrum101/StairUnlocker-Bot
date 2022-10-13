package utils

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
	"github.com/panjf2000/ants/v2"
)

type geoIP struct {
	Country string `json:"country"`
	Isp     string `json:"isp"`
	Query   string `json:"query"`
	Status  string `json:"status"`
}

const ipAPI = "http://ip-api.com/json?fields=25089"
const ipAPIBatch = "http://ip-api.com/batch?fields=25089"

func endIPTest(p C.Proxy) (ipInfo geoIP, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, err := urlToMetadata(ipAPI)
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

func GetIPList(proxiesList []C.Proxy, n int) ([]string, []string) {
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

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Dreamacro/clash/common/batch"
	C "github.com/Dreamacro/clash/constant"
)

type geoIP struct {
	Organization string `json:"organization"`
	Isp          string `json:"isp"`
	Country      string `json:"country"`
	Ip           string `json:"ip"`
	/*
		Longitude       float64 `json:"longitude"`
		Timezone        string  `json:"timezone"`
		Offset          int     `json:"offset"`
		Asn             int     `json:"asn"`
		AsnOrganization string  `json:"asn_organization"`
		Latitude        float64 `json:"latitude"`
		ContinentCode   string  `json:"continent_code"`
		CountryCode     string  `json:"country_code"`
	*/
}

func endIPTest(p C.Proxy, url string) (ipInfo geoIP, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addr, err := urlToMetadata(url)
	if err != nil {
		return
	}

	instance, err := p.DialContext(ctx, &addr)
	if err != nil {
		return
	}
	defer func(instance C.Conn) {
		err := instance.Close()
		if err != nil {
			return
		}
	}(instance)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "curl")
	req = req.WithContext(ctx)

	transport := &http.Transport{
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			return instance, nil
		},
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	buf, _ := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return
	}
	_ = json.Unmarshal(buf, &ipInfo)
	return
}

func entryIPTest(ip string) (ipInfo geoIP, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := "https://api.ip.sb/geoip/" + ip
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "curl")
	req = req.WithContext(ctx)

	client := http.Client{}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	buf, _ := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return
	}
	_ = json.Unmarshal(buf, &ipInfo)
	return
}

func entryIP(addr string) []string {
	li, _ := net.LookupHost(strings.Split(addr, ":")[0])
	return li
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
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(n))
	var endIPList []string
	var entryIPList []string
	for i := range proxiesList {
		p := proxiesList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			resp, err := endIPTest(p, "https://api.ip.sb/geoip")
			if err != nil {
				return nil, nil
			}
			endIPList = append(endIPList, fmt.Sprintf("%s - %s, ISP: %s", resp.Ip, resp.Country, resp.Isp))
			for i := range entryIP(p.Addr()) {
				resp, err = entryIPTest(entryIP(p.Addr())[i])
				if err != nil {
					return nil, nil
				}
				entryIPList = append(entryIPList, fmt.Sprintf("%s - %s, ISP: %s", resp.Ip, resp.Country, resp.Isp))
			}
			return nil, nil
		})
	}
	b.Wait()
	return deDuplication(entryIPList), deDuplication(endIPList)
}

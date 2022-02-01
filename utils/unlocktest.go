package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/go-resty/resty/v2"
)

func unlockTest(p *CheckAdapter) (t string, res bool, err error) {
	start := time.Now()
	resp, err := getURLResp(&p.Proxy, p.CheckURL)
	if err != nil {
		return
	}
	t = fmt.Sprintf("%dms", time.Since(start)/time.Millisecond)
	res = isUnlock(resp, p.CheckName, &p.Proxy)
	return
}

func urlToMetadata(rawURL string) (addr C.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}

	addr = C.Metadata{
		AddrType: C.AtypDomainName,
		Host:     u.Hostname(),
		DstIP:    nil,
		DstPort:  port,
	}
	return
}

func getURLResp(p *C.Proxy, url string) (resp *resty.Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return
	}

	instance, err := (*p).DialContext(ctx, &addr)
	if err != nil {
		return
	}

	defer func(instance C.Conn) {
		err := instance.Close()
		if err != nil {
			return
		}
	}(instance)

	resp, err = resty.New().SetTransport(&http.Transport{
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			return instance, err
		}}).SetCloseConnection(true).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62").
		R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	return resp, err
}

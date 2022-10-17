package provider

import (
	"context"
	"net"
	"net/http"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/go-resty/resty/v2"

	"github.com/thank243/StairUnlocker-Bot/model"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

type AbsStream interface {
	IsUnlock(p *C.Proxy) (model.StreamData, error)
}

type unlockProvider interface {
	create() AbsStream
}

func NewStreamList() []AbsStream {
	return []AbsStream{
		unlockProvider(new(netflix)).create(),
		unlockProvider(new(hbo)).create(),
		unlockProvider(new(youtube)).create(),
		unlockProvider(new(disney)).create(),
		unlockProvider(new(tvb)).create(),
		unlockProvider(new(abema)).create(),
		unlockProvider(new(bahamut)).create(),
		unlockProvider(new(connectivity)).create(),
	}
}

func getURLResp(p *C.Proxy, url string) (resp *resty.Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, err := utils.UrlToMetadata(url)
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

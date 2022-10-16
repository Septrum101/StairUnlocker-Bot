package utils

import (
	"fmt"
	"math"
	"net/url"
	"time"

	C "github.com/Dreamacro/clash/constant"
)

func UrlToMetadata(rawURL string) (addr C.Metadata, err error) {
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

func FormatTime(t time.Duration) string {
	timeStr := t.Round(time.Second).String()
	if t-24*time.Hour > 0 {
		day := int(math.Floor(float64(t / (24 * time.Hour))))
		timeStr = (t - time.Duration(day*24)*time.Hour).Round(time.Second).String()
		timeStr = fmt.Sprintf("%dd%s", day, timeStr)
	}
	return timeStr
}

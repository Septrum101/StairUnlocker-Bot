package utils

import (
	"fmt"
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
		NetWork: C.TCP,
		Host:    u.Hostname(),
		DstIP:   nil,
		DstPort: port,
	}
	return
}

func FormatTime(t time.Duration) string {
	d := t / (time.Hour * 24)
	if d > 0 {
		return fmt.Sprintf("%d days, %v", d, (t % (time.Hour * 24)).Round(time.Second))
	}
	return t.Round(time.Second).String()
}

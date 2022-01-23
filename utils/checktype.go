package utils

import (
	"strings"

	C "github.com/Dreamacro/clash/constant"
	"github.com/go-resty/resty/v2"
)

func isUnlock(r *resty.Response, testName string, p C.Proxy) bool {
	switch testName {
	case "Netflix", "HBO", "Youtube Premium":
		if r.StatusCode() < 300 {
			return true
		}
	case "Disney Plus":
		resp, err := getURLResp(p, "https://global.edge.bamgrid.com/token")
		if err != nil {
			return false
		}
		if r.StatusCode() < 300 && resp.StatusCode() != 403 {
			return true
		}
	case "TVB":
		if strings.Contains(string(r.Body()), "HK") {
			return true
		}
	case "Abema":
		if strings.Contains(string(r.Body()), "Country") {
			return true
		}
	case "Bahamut":
		if strings.Contains(string(r.Body()), "animeSn") {
			return true
		}
	}
	return false
}

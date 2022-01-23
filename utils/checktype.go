package utils

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func isUnlock(r *http.Response, testName string, req *http.Request) bool {
	switch testName {
	case "Netflix", "HBO", "Youtube Premium":
		if r.StatusCode < 300 {
			return true
		}
	case "Disney Plus":
		c := http.Client{}
		req.Host = "global.edge.bamgrid.com"
		req.URL, _ = url.Parse("https://global.edge.bamgrid.com/token")
		resp, _ := c.Do(req)
		if r.StatusCode < 300 && resp.StatusCode != 403 {
			return true
		}
	case "TVB":
		ctx, _ := io.ReadAll(r.Body)
		if strings.Contains(string(ctx), "HK") {
			return true
		}
	case "Abema":
		ctx, _ := io.ReadAll(r.Body)
		if strings.Contains(string(ctx), "Country") {
			return true
		}
	case "Bahamut":
		ctx, _ := io.ReadAll(r.Body)
		if strings.Contains(string(ctx), "animeSn") {
			return true
		}
	}
	return false
}

package utils

import (
	"io"
	"net/http"
	"strings"
)

func testResult(r *http.Response, testName string) bool {
	switch testName {
	case "Netflix", "HBO", "Disney Plus", "Youtube Premium":
		if r.StatusCode < 300 {
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

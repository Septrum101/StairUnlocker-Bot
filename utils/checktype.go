package utils

import "net/http"

func testResult(r *http.Response, testType int) bool {
	switch testType {
	// simple return statusCode
	case 0, 1, 2, 3:
		if r.StatusCode < 300 {
			return true
		}
	}
	return false
}

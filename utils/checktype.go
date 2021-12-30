package utils

import "net/http"

func testResult(r *http.Response, testType int) bool {
	switch {
	// simple return statusCode
	case testType == 0 || testType == 1 || testType == 2 || testType == 3:
		if r.StatusCode < 300 {
			return true
		}
	}
	return false
}

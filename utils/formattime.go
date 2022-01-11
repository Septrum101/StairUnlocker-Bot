package utils

import (
	"fmt"
	"math"
	"time"
)

func FormatTime(t time.Duration) string {
	timeStr := t.Round(time.Second).String()
	if t-24*time.Hour > 0 {
		day := int(math.Floor(float64(t / (24 * time.Hour))))
		timeStr = (t - time.Duration(day*24)*time.Hour).Round(time.Second).String()
		timeStr = fmt.Sprintf("%dd%s", day, timeStr)
	}
	return timeStr
}

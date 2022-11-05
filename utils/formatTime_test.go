package utils

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	tt := time.Second * 59
	t.Log(FormatTime(tt))
	tt += time.Minute * 59
	t.Log(FormatTime(tt))
	tt += time.Hour * 23
	t.Log(FormatTime(tt))
	tt += time.Hour * 24 * 365
	t.Log(FormatTime(tt))
}

package utils

import C "github.com/Dreamacro/clash/constant"

type checkParams struct {
	C.Proxy
	testName string
	testType int
	testURL  string
}

type CheckData struct {
	ProxyName   string
	StreamMedia string
	Latency     string
}

func getTestParams() []checkParams {
	return []checkParams{
		{testName: "Netflix", testURL: "https://www.netflix.com/title/70143836"},
		{testName: "HBO", testURL: "https://www.hbomax.com"},
		{testName: "Disney Plus", testURL: "https://www.disneyplus.com"},
		{testName: "Youtube Premium", testURL: "https://music.youtube.com"},
	}
}

package utils

import C "github.com/Dreamacro/clash/constant"

type CheckAdapter struct {
	C.Proxy
	CheckName string
	CheckURL  string
}

type CheckData struct {
	ProxyName   string
	StreamMedia string
	Latency     string
}

func GetCheckParams() []*CheckAdapter {
	return []*CheckAdapter{
		{CheckName: "Netflix", CheckURL: "https://www.netflix.com/title/70143836"},
		{CheckName: "HBO", CheckURL: "https://www.hbomax.com"},
		{CheckName: "Disney Plus", CheckURL: "https://www.disneyplus.com"},
		{CheckName: "Youtube Premium", CheckURL: "https://music.youtube.com"},
		{CheckName: "TVB", CheckURL: "https://www.mytvsuper.com/iptest.php"},
		{CheckName: "Abema", CheckURL: "https://api.abema.io/v1/ip/check?device=android"},
		{CheckName: "Bahamut", CheckURL: "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667"},
	}
}

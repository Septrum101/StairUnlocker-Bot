package user

import (
	"bytes"
	"github.com/Dreamacro/clash/log"
	"github.com/fogleman/gg"
	"github.com/thank243/StairUnlocker-Bot/config"
	"image/png"
	"sort"
	"strings"
	"time"
)

const W = 1000

func getBeginFix(i int) float64 {
	f := 2.3
	return W/f + (W-40-(W-40)/f)/4*float64(i)
}

func getStrWidth(dc *gg.Context, str string) (reStr string, width float64) {
	width, _ = dc.MeasureString(str)
	return str, width
}

func generatePNG(streamMediaUnlockMap map[string][]string) (*bytes.Buffer, error) {
	streamMediaNames := []string{"Netflix", "HBO", "Disney Plus", "Youtube Premium"}
	H := len(streamMediaUnlockMap)*25 + 100
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	// load font.
	path := strings.Split(config.ConfPath, "/")
	path[len(path)-1] = "msyh.ttc"
	err := dc.LoadFontFace(strings.Join(path, "/"), 15)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}

	str, strWidth := getStrWidth(dc, "StairUnlocker Bot")
	dc.DrawString(str, (W-strWidth)/2, 18)
	dc.SetLineWidth(0.4)
	// draw horizon lines
	for i := 0; i < len(streamMediaUnlockMap)+2; i++ {
		dc.DrawLine(20, 25+float64(i)*25, W-20, 25+float64(i)*25)
		if i == 0 {
			continue
		}
		dc.Stroke()
	}
	// draw vertical lines
	for i := 0; i < 2; i++ {
		dc.DrawLine(20+(W-40)*float64(i), 25, 20+(W-40)*float64(i), float64(H-50))
		dc.Stroke()
	}
	for i := 0; i < len(streamMediaNames); i++ {
		dc.DrawLine(getBeginFix(i), 25, getBeginFix(i), float64(H-50))
		dc.Stroke()
	}
	// header
	str, strWidth = getStrWidth(dc, "Node Name")
	dc.DrawString(str, 20+(getBeginFix(0)-20-strWidth)/2, 42.5)
	for i := range streamMediaNames {
		str, strWidth = getStrWidth(dc, streamMediaNames[i])
		dc.DrawString(str, getBeginFix(i)+(getBeginFix(1)-getBeginFix(0)-strWidth)/2, 42.5)
	}
	// context
	// sort nodes name
	var nameSort []string
	for i := range streamMediaUnlockMap {
		nameSort = append(nameSort, i)
	}
	sort.Strings(nameSort)

	n := 0
	for i := range nameSort {
		dc.DrawString(nameSort[i], 22, 67.5+float64(n)*25)
		for idx := range streamMediaUnlockMap[nameSort[i]] {
			if streamMediaUnlockMap[nameSort[i]][idx] != "" {
				dc.DrawString(streamMediaUnlockMap[nameSort[i]][idx]+"ms", 5+getBeginFix(idx), 67.5+float64(n)*25)
			} else {
				dc.SetRGB(1, 0, 0)
				str, strWidth = getStrWidth(dc, "None")
				dc.DrawString(str, getBeginFix(idx)+(getBeginFix(1)-getBeginFix(0)-strWidth)/2, 67.5+float64(n)*25)
				dc.SetRGB(0, 0, 0)
			}
		}
		n++
	}
	// bottom
	dc.DrawString("Timestamp: "+time.Now().Round(1*time.Second).String(), 20, float64(H)-30)
	str, strWidth = getStrWidth(dc, "Powered by @stairunlock_test_bot (https://git.io/Jyl5l)")
	dc.DrawString(str, W-20-strWidth, float64(H)-10)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	return buf, err
}

package user

import (
	"bytes"
	"fmt"
	"image/png"
	"sort"
	"strings"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/fogleman/gg"

	"github.com/thank243/StairUnlocker-Bot/config"
	"github.com/thank243/StairUnlocker-Bot/utils"
)

const W = 1500

func getBeginFix(i int, n int) float64 {
	f := 3.5
	return W/f + (W-40-(W-40)/f)/float64(n)*float64(i)
}

func getStrWidth(dc *gg.Context, str string) (reStr string, width float64) {
	width, _ = dc.MeasureString(str)
	return str, width
}

func generatePNG(streamDataList []utils.StreamData, streamMediaNames []string) (*bytes.Buffer, error) {
	// merge node unlock data
	nodeMap := make(map[string][]utils.StreamData)
	for i := range streamDataList {
		nodeMap[streamDataList[i].ProxyName] = append(nodeMap[streamDataList[i].ProxyName], streamDataList[i])
	}
	var sortedNode [][]utils.StreamData

	for idx := range nodeMap {
		sort.Slice(nodeMap[idx], func(i, j int) bool {
			return nodeMap[idx][i].Name < nodeMap[idx][j].Name
		})

		sortedNode = append(sortedNode, nodeMap[idx])
	}

	sort.Slice(sortedNode, func(i, j int) bool {
		return sortedNode[i][0].ProxyName < sortedNode[j][0].ProxyName
	})

	H := len(sortedNode)*25 + 110
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

	str, strWidth := getStrWidth(dc, fmt.Sprintf("StairUnlocker Bot %s", C.Version))
	dc.DrawString(str, (W-strWidth)/2, 20)
	dc.SetLineWidth(0.4)
	// draw horizon lines
	for i := 0; i < len(sortedNode)+2; i++ {
		dc.DrawLine(20, 35+float64(i)*25, W-20, 35+float64(i)*25)
		if i == 0 {
			continue
		}
		dc.Stroke()
	}
	// draw vertical lines
	for i := 0; i < 2; i++ {
		dc.DrawLine(20+(W-40)*float64(i), 35, 20+(W-40)*float64(i), float64(H-50))
		dc.Stroke()
	}
	for i := 0; i < len(streamMediaNames); i++ {
		dc.DrawLine(getBeginFix(i, len(streamMediaNames)), 35, getBeginFix(i, len(streamMediaNames)), float64(H-50))
		dc.Stroke()
	}
	// header
	str, strWidth = getStrWidth(dc, "Node Name")
	dc.DrawString(str, 20+(getBeginFix(0, len(streamMediaNames))-20-strWidth)/2, 52.5)
	for i := range streamMediaNames {
		str, strWidth = getStrWidth(dc, streamMediaNames[i])
		dc.DrawString(str, getBeginFix(i, len(streamMediaNames))+(getBeginFix(1, len(streamMediaNames))-getBeginFix(0, len(streamMediaNames))-strWidth)/2, 52.5)
	}

	// context
	n := 0
	for i := range sortedNode {
		dc.DrawString(sortedNode[i][0].ProxyName, 22, 77.5+float64(n)*25)
		for idx := range sortedNode[i] {
			if sortedNode[i][idx].Unlock {
				dc.DrawString(sortedNode[i][idx].Latency, 5+getBeginFix(idx, len(streamMediaNames)), 77.5+float64(n)*25)
			} else {
				dc.SetRGB(1, 0, 0)
				str, strWidth = getStrWidth(dc, "None")
				dc.DrawString(str, getBeginFix(idx, len(streamMediaNames))+(getBeginFix(1, len(streamMediaNames))-getBeginFix(0, len(streamMediaNames))-strWidth)/2, 77.5+float64(n)*25)
				dc.SetRGB(0, 0, 0)
			}
		}
		n++
	}

	// bottom
	dc.DrawString(fmt.Sprintf("Timestamp: %s", time.Now().UTC().Format(time.RFC3339)), 20, float64(H)-25)
	str, strWidth = getStrWidth(dc, "Powered by @stairunlock_test_bot (https://git.io/Jyl5l)")
	dc.DrawString(str, W-20-strWidth, float64(H)-10)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	return buf, err
}

package utils

import (
	"context"
	"github.com/Dreamacro/clash/common/batch"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"strconv"
)

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (streamMediaUnlockList []CheckData) {
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(n))
	// counts buffer channel
	ch := make(chan int, 16)
	defer close(ch)
	testParams := getTestParams()
	var wrapList []checkParams
	for i := range proxiesList {
		for idx := range testParams {
			wrapList = append(wrapList, checkParams{
				Proxy:    proxiesList[i],
				testName: testParams[idx].testName,
				testType: idx,
				testURL:  testParams[idx].testURL,
			})
		}
	}
	curr, total := 0, len(wrapList)
	for i := range wrapList {
		p := wrapList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			latency, resp, err := streamMediaUnlockTest(p)
			if err != nil {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s : %s", curr, total, p.Name(), err.Error())
			} else if resp {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s | %s Unlock", curr, total, p.Name(), p.testName)
				streamMediaUnlockList = append(streamMediaUnlockList, CheckData{p.Name(), p.testName, strconv.Itoa(int(latency))})
			} else {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s | None", curr, total, p.Name())
			}
			return nil, nil
		})
	}
	b.Wait()
	return
}

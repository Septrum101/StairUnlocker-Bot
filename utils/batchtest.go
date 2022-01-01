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
	testParams := GetCheckParams()
	var wrapList []CheckAdapter
	for i := range proxiesList {
		for idx := range testParams {
			wrapList = append(wrapList, CheckAdapter{
				Proxy:     proxiesList[i],
				CheckName: testParams[idx].CheckName,
				CheckURL:  testParams[idx].CheckURL,
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
				log.Debugln("(%d/%d) %s | %s Unlock", curr, total, p.Name(), p.CheckName)
				streamMediaUnlockList = append(streamMediaUnlockList, CheckData{p.Name(), p.CheckName, strconv.Itoa(int(latency))})
			} else {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s | %s None", curr, total, p.Name(), p.CheckName)
			}
			return nil, nil
		})
	}
	b.Wait()
	return
}

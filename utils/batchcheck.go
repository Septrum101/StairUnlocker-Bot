package utils

import (
	"sync"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/panjf2000/ants/v2"
)

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (streamMediaUnlockList []CheckData) {
	var (
		wg       sync.WaitGroup
		wrapList []CheckAdapter
	)
	// wrap node list for test.
	for i := range proxiesList {
		for idx := range GetCheckParams() {
			wrapList = append(wrapList, CheckAdapter{
				Proxy:     proxiesList[i],
				CheckName: GetCheckParams()[idx].CheckName,
				CheckURL:  GetCheckParams()[idx].CheckURL,
			})
		}
	}
	// prefix for node name on log.
	curr, total := 0, len(wrapList)
	//initial pool
	pool, err := ants.NewPoolWithFunc(n, func(i interface{}) {
		p := i.(CheckAdapter)
		latency, resp, err := unlockTest(&p)
		if err != nil {
			curr++
			log.Debugln("(%d/%d) %s : %s", curr, total, p.Name(), err.Error())
		} else if resp {
			curr++
			log.Debugln("(%d/%d) %s | %s Unlock", curr, total, p.Name(), p.CheckName)
			streamMediaUnlockList = append(streamMediaUnlockList, CheckData{p.Name(), p.CheckName, latency})
		} else {
			curr++
			log.Debugln("(%d/%d) %s | %s None", curr, total, p.Name(), p.CheckName)
		}
		wg.Done()
	})
	defer pool.Release()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	for i := range wrapList {
		wg.Add(1)
		err = pool.Invoke(wrapList[i])
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	wg.Wait()
	return
}

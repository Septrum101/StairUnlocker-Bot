package utils

import (
	"sync"
	"sync/atomic"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/panjf2000/ants/v2"

	"github.com/thank243/StairUnlocker-Bot/model"
)

var l sync.RWMutex

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (streamDataList []model.StreamData) {
	type combineProxy struct {
		proxy  C.Proxy
		stream absStream
	}

	var (
		wg sync.WaitGroup
		cp []combineProxy
	)
	streamList := NewStreamList()

	for i := range proxiesList {
		for ii := range streamList {
			cp = append(cp, combineProxy{
				proxy:  proxiesList[i],
				stream: streamList[ii],
			})
		}
	}
	// prefix for node name on log.
	curr, total := int32(0), len(cp)
	//initial pool
	pool, err := ants.NewPoolWithFunc(n, func(i interface{}) {
		c := i.(combineProxy)
		result, err := c.stream.isUnlock(&c.proxy)
		atomic.AddInt32(&curr, 1)
		if err != nil {
			log.Debugln("(%d/%d) %s : %s", atomic.LoadInt32(&curr), total, c.proxy.Name(), err.Error())
		} else if result.Unlock {
			log.Debugln("(%d/%d) %s | %s Unlock", atomic.LoadInt32(&curr), total, c.proxy.Name(), result.Name)
		} else {
			log.Debugln("(%d/%d) %s | %s None", atomic.LoadInt32(&curr), total, c.proxy.Name(), result.Name)
		}
		l.Lock()
		streamDataList = append(streamDataList, result)
		l.Unlock()
		wg.Done()
	})
	defer pool.Release()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	for i := range cp {
		wg.Add(1)
		err = pool.Invoke(cp[i])
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	wg.Wait()
	return
}

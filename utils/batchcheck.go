package utils

import (
	"sync"
	"sync/atomic"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/panjf2000/ants/v2"
)

var l sync.RWMutex

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (streamDataList []StreamData) {
	var (
		wg sync.WaitGroup
	)
	streamList := []absStream{
		unlockFactory(new(netflix)).create(),
		unlockFactory(new(hbo)).create(),
		unlockFactory(new(youtube)).create(),
		unlockFactory(new(disney)).create(),
		unlockFactory(new(tvb)).create(),
		unlockFactory(new(abema)).create(),
		unlockFactory(new(bahamut)).create(),
	}
	// prefix for node name on log.
	curr, total := int32(0), len(proxiesList)*len(streamList)
	//initial pool
	pool, err := ants.NewPoolWithFunc(n, func(i interface{}) {
		p := i.(C.Proxy)
		for ii := range streamList {
			wg.Add(1)
			j := ii
			go func() {
				result, err := streamList[j].isUnlock(&p)
				if err != nil {
					atomic.AddInt32(&curr, 1)
					log.Debugln("(%d/%d) %s : %s", curr, total, p.Name(), err.Error())
				} else if result.Unlock {
					atomic.AddInt32(&curr, 1)
					log.Debugln("(%d/%d) %s | %s Unlock", curr, total, p.Name(), result.Name)
				} else {
					atomic.AddInt32(&curr, 1)
					log.Debugln("(%d/%d) %s | %s None", curr, total, p.Name(), result.Name)
				}
				l.Lock()
				streamDataList = append(streamDataList, result)
				l.Unlock()
				wg.Done()
			}()
		}
		wg.Done()
	})
	defer pool.Release()
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	for i := range proxiesList {
		wg.Add(1)
		err = pool.Invoke(proxiesList[i])
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	wg.Wait()
	return
}

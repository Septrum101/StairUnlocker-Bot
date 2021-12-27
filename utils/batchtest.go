package utils

import (
	"context"
	"github.com/Dreamacro/clash/common/batch"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

func getStreamMedias() []string {
	return []string{
		"https://www.netflix.com/title/70143836",
		"https://www.hbomax.com",
		"https://www.disneyplus.com",
		"https://music.youtube.com",
	}
}

func getSteamMediaNames() []string {
	return []string{"Netflix", "HBO", "DisneyPlus", "Youtube Premium"}
}

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (streamMediaList []string, latencyMap map[string][]uint16) {
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(n))
	// counts buffer channel
	ch := make(chan int, 16)
	latencyMap = make(map[string][]uint16)
	lockMap := syncMap{Map: latencyMap}
	defer close(ch)
	curr, total := 0, len(proxiesList)*len(getStreamMedias())
	for i := range proxiesList {
		p := proxiesList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			for idx := range getStreamMedias() {
				latency, sCode, err := streamMediaUnlockTest(p, getStreamMedias()[idx])
				if err != nil {
					ch <- 1
					curr += <-ch
					log.Debugln("(%d/%d) %s : %s", curr, total, p.Name(), err.Error())
					lockMap.Lock()
					latencyMap[p.Name()] = append(latencyMap[p.Name()], 0)
					lockMap.Unlock()
				} else if sCode < 300 {
					ch <- 1
					curr += <-ch
					log.Debugln("(%d/%d) %s | %s Unlock", curr, total, p.Name(), getSteamMediaNames()[idx])
					lockMap.Lock()
					latencyMap[p.Name()] = append(latencyMap[p.Name()], latency)
					lockMap.Unlock()
					streamMediaList = append(streamMediaList, p.Name())
				} else {
					ch <- 1
					curr += <-ch
					log.Debugln("(%d/%d) %s | None", curr, total, p.Name())
					lockMap.Lock()
					latencyMap[p.Name()] = append(latencyMap[p.Name()], 0)
					lockMap.Unlock()
				}
			}
			return nil, nil
		})
	}
	b.Wait()
	return
}

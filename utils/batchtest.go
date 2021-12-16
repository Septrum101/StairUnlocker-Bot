package utils

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/common/batch"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (netflixList []string, latencyList []string) {
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(n))
	// counts buffer channel
	ch := make(chan int, 16)
	// close channel
	defer close(ch)
	curr, total := 0, len(proxiesList)
	for i := range proxiesList {
		p := proxiesList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			latency, sCode, err := netflixTest(p, "https://www.netflix.com/title/70143836")
			if err != nil {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s : %s", curr, total, p.Name(), err.Error())
			} else if sCode == 200 {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s : latency = %v ms | Full Unlock", curr, total, p.Name(), latency)
				latencyList = append(latencyList, fmt.Sprintf("%s - %dms", p.Name(), latency))
				netflixList = append(netflixList, p.Name())
			} else {
				ch <- 1
				curr += <-ch
				log.Debugln("(%d/%d) %s : latency = %v ms | None", curr, total, p.Name(), latency)
			}
			return nil, nil
		})
	}
	b.Wait()
	return
}

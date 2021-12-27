package utils

import "sync"

type syncMap struct {
	sync.RWMutex
	Map map[string][]uint16
}

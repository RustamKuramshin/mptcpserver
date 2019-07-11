package lib

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	value      uint32
	lastAccess int64
}

type TTLMap struct {
	m map[string]*item
	l sync.Mutex
}

func NewTTLMap(ln int, maxTTL int, printTime int) (ttlmap *TTLMap) {
	ttlmap = &TTLMap{m: make(map[string]*item, ln)}
	go func() {
		for now := range time.Tick(time.Second) {
			ttlmap.l.Lock()
			for k, v := range ttlmap.m {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(ttlmap.m, k)
				}
			}
			ttlmap.l.Unlock()
		}
	}()

	go func() {
		for range time.Tick(time.Duration(printTime) * time.Second) {
			ttlmap.l.Lock()
			fmt.Printf("TTLMap content... \n")
			if ttlmap.Len() > 0 {
				for k, v := range ttlmap.m {
					fmt.Printf("key: %v | value: %v \n", k, v.value)
				}
			} else {
				fmt.Printf("TTLMap is empty \n")
			}
			ttlmap.l.Unlock()
		}
	}()

	return
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k string, v uint32) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &item{value: v}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap) Get(k string) (v uint32) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.value
		it.lastAccess = time.Now().Unix()
	}
	m.l.Unlock()
	return

}

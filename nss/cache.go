package main

import (
	"sync"
	"time"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/settings"
)

var store map[string]*cacheObject

type cacheObject struct {
	userGroup *stns.Attributes
	createAt  time.Time
	err       error
}

func readCache(path string) (stns.Attributes, error) {
	m := sync.RWMutex{}
	m.RLock()
	defer m.RUnlock()

	initCache()

	c, exist := store[path]
	if exist {
		// cache expire 10 minute
		if time.Now().Sub(c.createAt) > time.Minute*settings.CACHE_TIME {
			delete(store, path)
			return nil, nil
		} else if c.err != nil {
			return nil, c.err
		} else {
			return *c.userGroup, c.err
		}
	}
	return nil, nil
}

func writeCache(path string, attr stns.Attributes, err error) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	initCache()

	store[path] = &cacheObject{&attr, time.Now(), err}
}

func initCache() {
	if len(store) == 0 {
		store = map[string]*cacheObject{}
	}
}

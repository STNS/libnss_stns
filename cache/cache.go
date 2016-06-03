package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func Read(path string) (stns.Attributes, error) {
	m := sync.RWMutex{}
	m.RLock()
	defer m.RUnlock()

	Init()

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

func Write(path string, attr stns.Attributes, err error) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	Init()

	store[path] = &cacheObject{&attr, time.Now(), err}
}

func SaveResultList(resourceType string, list stns.Attributes) {
	j, err := json.Marshal(list)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(settings.CACHE_DIR+"/.stns_"+resourceType+"_cache", j, 0644)
}

func LastResultList(resourceType string) *stns.Attributes {
	var attr stns.Attributes
	f, err := ioutil.ReadFile(settings.CACHE_DIR + "/.stns_" + resourceType + "_cache")
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(f, &attr)
	if err != nil {
		log.Println(err)
	}
	return &attr
}
func Init() {
	if len(store) == 0 {
		store = map[string]*cacheObject{}
	}
}

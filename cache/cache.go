package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/settings"
	_gocache "github.com/patrickmn/go-cache"
)

var store *_gocache.Cache

type cacheObject struct {
	userGroup *stns.Attributes
	createAt  time.Time
	err       error
}

func Read(path string) (stns.Attributes, error) {
	Init()
	c, exist := store.Get(path)
	if exist {
		co := c.(*cacheObject)
		if co.err != nil {
			return nil, co.err
		} else {
			return *co.userGroup, co.err
		}
	}
	return nil, nil
}

func Write(path string, attr stns.Attributes, err error) {
	Init()
	store.Set(path, &cacheObject{&attr, time.Now(), err}, _gocache.DefaultExpiration)
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
	if store == nil {
		store = _gocache.New(time.Minute*settings.CACHE_TIME, 60*time.Second)
	}
}

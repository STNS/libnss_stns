package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/settings"
	_gocache "github.com/patrickmn/go-cache"
)

var store = _gocache.New(time.Minute*settings.CACHE_TIME, 60*time.Second)
var workDir = settings.WORK_DIR

type cacheObject struct {
	userGroup *stns.Attributes
	createAt  time.Time
	err       error
}

func SetWorkDir(path string) {
	workDir = path
}

func Read(path string) (stns.Attributes, error) {
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
	store.Set(path, &cacheObject{&attr, time.Now(), err}, _gocache.DefaultExpiration)
}

func SaveResultList(resourceType string, list stns.Attributes) {
	j, err := json.Marshal(list)
	if err != nil {
		log.Println(err)
		return
	}

	if err := os.MkdirAll(workDir, 0777); err != nil {
		log.Println(err)
	}
	ioutil.WriteFile(workDir+"/.libnss_stns_"+resourceType+"_cache", j, 0777)
}

func LastResultList(resourceType string) *stns.Attributes {
	var attr stns.Attributes
	f, err := ioutil.ReadFile(workDir + "/.libnss_stns_" + resourceType + "_cache")
	if err != nil {
		log.Println(err)
		return &attr
	}

	err = json.Unmarshal(f, &attr)
	if err != nil {
		log.Println(err)
	}
	return &attr
}

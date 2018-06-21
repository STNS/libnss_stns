package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/settings"
	_gocache "github.com/pyama86/go-cache"
)

var attrStore = _gocache.New(time.Second*settings.CACHE_TIME, 0*time.Second)
var lockStore = _gocache.New(time.Second*settings.LOCK_TIME, 0*time.Second)

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
	c, exist := attrStore.Get(path)
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
	attrStore.Set(path, &cacheObject{&attr, time.Now(), err}, _gocache.DefaultExpiration)
}

func ReadMinID(resourceType string) int {
	return readID(resourceType, "min")
}

func ReadMaxID(resourceType string) int {
	return readID(resourceType, "max")
}

func readID(resourceType, minMax string) int {
	n, exist := attrStore.Get(resourceType + "_" + minMax + "_id")
	if exist {
		id := n.(int)
		return id
	}
	return 0
}

func WriteID(resourceType, minMax string, id int) {
	attrStore.Set(resourceType+"_"+minMax+"_id", id, _gocache.DefaultExpiration)
}

func SaveResultList(resourceType string, list stns.Attributes) {
	j, err := json.Marshal(list)
	if err != nil {
		log.Println(err)
		return
	}

	if err := os.MkdirAll(workDir, 0770); err != nil {
		log.Println(err)
		return
	}
	f := workDir + "/.libnss_stns_" + resourceType + "_cache"

	if err := ioutil.WriteFile(f, j, os.ModePerm); err != nil {
		log.Println(err)
		return
	}

	n, err := user.LookupGroup("nscd")
	if err != nil {
		log.Println(err)
		return
	}
	gid, err := strconv.Atoi(n.Gid)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.Chown(f, 0, gid)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.Chmod(f, 0640)
	if err != nil {
		log.Println(err)
		return
	}
}

func LastResultList(resourceType string) *stns.Attributes {
	var attr stns.Attributes
	f := workDir + "/.libnss_stns_" + resourceType + "_cache"

	if _, err := os.Stat(f); err == nil {
		f, err := ioutil.ReadFile(f)
		if err != nil {
			log.Println(err)
			return &attr
		}

		err = json.Unmarshal(f, &attr)
		if err != nil {
			log.Println(err)
		}
	}
	return &attr
}

func LockEndPoint(path string) {
	lockStore.Set(path+"_lock", true, _gocache.DefaultExpiration)

	err := lockStore.SaveFile(settings.LOCK_FILE)
	if err != nil {
		log.Printf("lock file write error:%s", err.Error())
	}

	os.Chmod(settings.LOCK_FILE, 0777)
}

func IsLockEndPoint(path string) bool {
	_, e1 := lockStore.Get(path + "_lock")
	if e1 {
		return true
	} else {
		err := lockStore.LoadFile(settings.LOCK_FILE)
		if err != nil {
			return false
		}

		_, e2 := lockStore.Get(path + "_lock")
		if e2 {
			return true
		}
	}
	return false
}

func Flush() {
	attrStore.Flush()
	lockStore.Flush()
}

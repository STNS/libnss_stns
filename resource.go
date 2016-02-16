package main

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/logger"
)

const NSS_STATUS_TRYAGAIN = -2
const NSS_STATUS_SUCCESS = 1
const NSS_STATUS_NOTFOUND = 0

var ConfigFileName = "/etc/stns/libnss_stns.conf"
var Config *config.Config

type LinuxResource interface {
	setCStruct(attribute.UserGroups)
}

var Cache map[string]*CacheObject

type CacheObject struct {
	userGroup *attribute.UserGroups
	createAt  time.Time
	err       error
}

func get(paths ...string) (attribute.UserGroups, error) {
	logger.Setlog()
	path := strings.Join(paths, "/")

	u, err := readCache(path)
	if u != nil || err != nil {
		return u, err
	}

	if Config == nil {
		config, err := config.Load(ConfigFileName)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		Config = config
	}
	// default negative cache
	writeCache(path, nil, errors.New(path+" is not fond"))
	var attr attribute.UserGroups

	out, err := exec.Command(Config.WrapperCommand, path).Output()

	if err != nil || out == nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(out, &attr)
	if err != nil {
		return nil, err
	}
	writeCache(path, attr, nil)

	return attr, nil
}

func setResource(linux LinuxResource, resource_type, column string, value string) int {
	resource, err := get(resource_type, column, value)
	if err != nil {
		return NSS_STATUS_TRYAGAIN
	}

	if len(resource) > 0 {
		linux.setCStruct(resource)
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_NOTFOUND
}

func setNextResource(linux LinuxResource, list attribute.UserGroups, position *int) int {
	keys := keys(list)
	if *position != NSS_STATUS_TRYAGAIN && len(keys) > *position && keys[*position] != "" {
		name := keys[*position]
		resource := attribute.UserGroups{
			name: list[name],
		}

		linux.setCStruct(resource)
		*position++
		return NSS_STATUS_SUCCESS
	} else if *position == NSS_STATUS_TRYAGAIN {
		return NSS_STATUS_TRYAGAIN
	}
	return NSS_STATUS_NOTFOUND
}

func setList(resource_type string, list attribute.UserGroups, position *int) {
	// reset value
	resetList(list, position)

	resource, err := get(resource_type, "list", "")
	if err != nil {
		*position = NSS_STATUS_TRYAGAIN
		return
	}

	if len(resource) > 0 {
		for k, v := range resource {
			list[k] = v
		}
	}
}

func resetList(list attribute.UserGroups, position *int) {
	// reset value
	*position = 0
	for k, _ := range list {
		delete(list, k)
	}
}

func keys(list attribute.UserGroups) []string {
	ks := []string{}
	for k, _ := range list {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

func readCache(path string) (attribute.UserGroups, error) {
	m := sync.RWMutex{}
	m.RLock()
	defer m.RUnlock()

	if len(Cache) == 0 {
		Cache = map[string]*CacheObject{}
	}

	c, exist := Cache[path]
	if exist {
		// cache expire 10 minute
		if time.Now().Sub(c.createAt) > time.Minute*10 {
			delete(Cache, path)
			return nil, nil
		} else if c.err != nil {
			return nil, c.err
		} else {
			return *c.userGroup, c.err
		}
	}
	return nil, nil
}

func writeCache(path string, attr attribute.UserGroups, err error) {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	if len(Cache) == 0 {
		Cache = map[string]*CacheObject{}
	}

	Cache[path] = &CacheObject{&attr, time.Now(), err}
}

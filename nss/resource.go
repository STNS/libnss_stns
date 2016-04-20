package main

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"sort"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/config"
)

const NSS_STATUS_TRYAGAIN = -2
const NSS_STATUS_SUCCESS = 1
const NSS_STATUS_NOTFOUND = 0

var conf *config.Config

type LinuxResource interface {
	setCStruct(stns.Attributes) int
}

func get(paths ...string) (stns.Attributes, error) {
	path := strings.Join(paths, "/")

	u, err := readCache(path)
	if u != nil || err != nil {
		return u, err
	}

	if conf == nil {
		c, err := config.Load("/etc/stns/libnss_stns.conf")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		conf = c
	}

	// deault negative cache
	writeCache(path, nil, errors.New(path+" is not fond"))
	out, _ := exec.Command(conf.WrapperCommand, path).Output()

	var attr stns.Attributes
	err = json.Unmarshal(out, &attr)
	if err != nil {
		return nil, err
	}
	writeCache(path, attr, nil)
	return attr, nil
}

func set(linux LinuxResource, resource_type, column string, value string) int {
	resource, err := get(resource_type, column, value)
	if err != nil {
		return NSS_STATUS_TRYAGAIN
	}

	if len(resource) > 0 {
		return linux.setCStruct(resource)
	}
	return NSS_STATUS_NOTFOUND
}

func setByList(linux LinuxResource, list stns.Attributes, position *int) int {
	keys := keys(list)
L:
	if *position != NSS_STATUS_TRYAGAIN && len(keys) > *position && keys[*position] != "" {
		name := keys[*position]
		resource := stns.Attributes{
			name: list[name],
		}

		*position++
		result := linux.setCStruct(resource)

		// lack of data
		if result == NSS_STATUS_NOTFOUND {
			goto L
		}

		return result
	} else if *position == NSS_STATUS_TRYAGAIN {
		return NSS_STATUS_TRYAGAIN
	}
	return NSS_STATUS_NOTFOUND
}

func initList(resource_type string, list stns.Attributes, position *int) int {
	// reset value
	purgeList(list, position)

	resource, err := get(resource_type, "list")
	if err != nil {
		*position = NSS_STATUS_TRYAGAIN
		return NSS_STATUS_TRYAGAIN
	}

	if len(resource) > 0 {
		for k, v := range resource {
			list[k] = v
		}
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_NOTFOUND
}

func purgeList(list stns.Attributes, position *int) {
	// reset value
	*position = 0
	for k, _ := range list {
		delete(list, k)
	}
}

func keys(list stns.Attributes) []string {
	ks := []string{}
	for k, _ := range list {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

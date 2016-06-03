package main

import (
	"errors"
	"log"
	"sort"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/cache"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/request"
)

import "C"

const (
	NSS_STATUS_NOTFOUND = C.int(0)
	NSS_STATUS_SUCCESS  = C.int(1)
	NSS_STATUS_UNAVAIL  = C.int(2)
	NSS_STATUS_TRYGAIN  = C.int(3)
	LIST_EMPTY          = -1
)

var conf *config.Config

type LinuxResource interface {
	setCStruct(stns.Attributes) C.int
}

func get(paths ...string) (stns.Attributes, error) {
	if conf == nil {
		c, err := config.Load("/etc/stns/libnss_stns.conf")
		if err != nil {
			log.Print(err)
			return nil, err
		}
		conf = c
	}

	r, err := request.NewRequest(conf, paths...)
	u, err := cache.Read(r.ApiPath)
	if u != nil || err != nil {
		return u, err
	}

	// deault negative cache
	cache.Write(r.ApiPath, nil, errors.New(r.ApiPath+" is not fond"))
	res, err := r.GetByWrapperCmd()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	cache.Write(r.ApiPath, *res.Items, nil)
	return *res.Items, nil
}

func set(linux LinuxResource, resourceType, column string, value string) C.int {
	resource, err := get(resourceType, column, value)
	if err != nil {
		return NSS_STATUS_UNAVAIL
	}

	if len(resource) > 0 {
		return linux.setCStruct(resource)
	}
	return NSS_STATUS_NOTFOUND
}

func setByList(linux LinuxResource, list stns.Attributes, position *int) C.int {
	keys := keys(list)
L:
	if *position != LIST_EMPTY && len(keys) > *position && keys[*position] != "" {
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
	} else if *position == LIST_EMPTY {
		return NSS_STATUS_UNAVAIL
	}
	return NSS_STATUS_NOTFOUND
}

func initList(resourceType string, list stns.Attributes, position *int) C.int {
	var attr stns.Attributes
	var err error
	// reset value
	purgeList(list, position)

	attr, err = get(resourceType, "list")

	if err != nil {
		// When the error refers to the last result.
		// This is supposed to when the server is restarted
		attr = *cache.LastResultList(resourceType)
		if attr != nil {
			goto C
		}

		*position = LIST_EMPTY
		return NSS_STATUS_UNAVAIL
	}
	cache.SaveResultList(resourceType, attr)
C:
	if len(attr) > 0 {
		for k, v := range attr {
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

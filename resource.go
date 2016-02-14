package main

import (
	"log"
	"sort"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/libnss_stns/request"
)

const NSS_STATUS_TRYAGAIN = -2
const NSS_STATUS_SUCCESS = 1
const NSS_STATUS_NOTFOUND = 0

type LinuxResource interface {
	setCStruct(attribute.UserGroups)
}

func get(resource_type, column, value string) (attribute.UserGroups, error) {
	req, err := request.NewRequest(resource_type, column, value)
	if err != nil {
		return nil, err
	}

	resource, err := req.Get()

	if err != nil {
		return nil, err
		log.Print(err)
	}
	return resource, nil
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

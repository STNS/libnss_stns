package main

import (
	"log"
	"sort"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/libnss_stns/request"
)

/*
#include <grp.h>
#include <pwd.h>
#include <shadow.h>
#include <sys/types.h>
*/
import "C"

const NSS_STATUS_TRYAGAIN = -2
const NSS_STATUS_SUCCESS = 1
const NSS_STATUS_NOTFOUND = 0

type Resource struct {
	resource_type string
}

type EntryResource struct {
	*Resource
	list     attribute.UserGroups
	position *int
}

type LinuxResource interface {
	setCStruct(attribute.UserGroups)
}

func (r *Resource) get(column string, value string) (attribute.UserGroups, error) {
	req, err := request.NewRequest(r.resource_type, column, value)
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

func (r *Resource) setResource(linux LinuxResource, column string, value string) int {
	resource, err := r.get(column, value)
	if err != nil {
		return NSS_STATUS_TRYAGAIN
	}

	if len(resource) > 0 {
		linux.setCStruct(resource)
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_NOTFOUND
}

func (e *EntryResource) setNextResource(linux LinuxResource) int {
	keys := e.keys()
	if *e.position != NSS_STATUS_TRYAGAIN && len(keys) > *e.position && keys[*e.position] != "" {
		name := keys[*e.position]
		resource := attribute.UserGroups{
			name: e.list[name],
		}

		linux.setCStruct(resource)
		*e.position++
		return NSS_STATUS_SUCCESS
	} else if *e.position == NSS_STATUS_TRYAGAIN {
		return NSS_STATUS_TRYAGAIN
	}
	return NSS_STATUS_NOTFOUND
}

func (e *EntryResource) setList() {
	// reset value
	e.resetList()

	resource, err := e.get("list", "")
	if err != nil {
		*e.position = NSS_STATUS_TRYAGAIN
	}

	if len(resource) > 0 {
		for k, v := range resource {
			e.list[k] = v
		}
	}
	return
}

func (e *EntryResource) resetList() {
	// reset value
	*e.position = 0
	for k, _ := range e.list {
		delete(e.list, k)
	}
}

func (e *EntryResource) keys() []string {
	ks := []string{}
	for k, _ := range e.list {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

package main

import (
	"log"
	"sort"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/request"
)

func setList(resource string, list map[string]*attribute.All, pos *int) {
	// reset value
	resetList(list, pos)

	l, err := request.GetList(resource)
	if err != nil {
		log.Print(err)
		return
	}

	for k, v := range l {
		list[k] = v
	}
}

func resetList(list map[string]*attribute.All, pos *int) {
	// reset value
	*pos = 0
	for k, _ := range list {
		delete(list, k)
	}
}

func getKeys(m map[string]*attribute.All) []string {
	ks := []string{}
	for k, _ := range m {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

func getNextResource(list map[string]*attribute.All, pos *int) map[string]*attribute.All {
	keys := getKeys(list)
	if len(keys) > *pos && keys[*pos] != "" {
		name := keys[*pos]
		resource := map[string]*attribute.All{
			name: list[name],
		}
		*pos++
		return resource
	}
	return nil
}

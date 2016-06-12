package libstns

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/cache"
)

// http://www.gnu.org/software/libc/manual/html_node/NSS-Modules-Interface.html
const (
	NSS_STATUS_SUCCESS  = 1
	NSS_STATUS_NOTFOUND = 0
	NSS_STATUS_UNAVAIL  = -1
	NSS_STATUS_TRYGAIN  = -2
)

const (
	NSS_LIST_PRESET = 0
	NSS_LIST_PURGE  = 1
)

type SetNss func(stns.Attributes, interface{}, interface{}) int

type Nss struct {
	config  *Config
	request *Request
	rtype   string
	column  string
	value   string
}

func NewNss(r, c, v string) (*Nss, error) {
	config, err := LoadConfig("/etc/stns/libnss_stns.conf")
	if err != nil {
		return nil, err
	}

	req, err := NewRequest(config, r, c, v)
	if err != nil {
		return nil, err
	}

	return &Nss{
		config:  config,
		request: req,
		rtype:   r,
		column:  c,
		value:   v,
	}, nil
}
func (n *Nss) Get() (stns.Attributes, error) {
	u, err := cache.Read(n.request.ApiPath)
	if u != nil || err != nil {
		return u, err
	}

	res, err := n.request.GetByWrapperCmd()

	if err != nil {
		return nil, err
	}

	if n.column == "id" {
		cache.WriteMinId(n.rtype, res.MetaData.MinId)
	}

	if res.Items == nil {
		return nil, fmt.Errorf("resource notfound %s/%s/%s", n.rtype, n.column, n.value)
	}

	cache.Write(n.request.ApiPath, *res.Items, nil)
	return *res.Items, nil
}

func (n *Nss) Set(s SetNss, entry, presult interface{}) int {
	id, _ := strconv.Atoi(n.column)
	if n.column != "id" || (n.column == "id" && cache.ReadMinId(n.rtype) <= id) {
		resource, err := n.Get()
		if err != nil {
			log.Print(err)
			return NSS_STATUS_UNAVAIL
		}

		if len(resource) > 0 {
			return s(resource, entry, presult)
		}
	}
	return NSS_STATUS_NOTFOUND
}

func (n *Nss) SetByList(s SetNss, entry, presult interface{}, list stns.Attributes, position *int) int {
	keys := n.keys(list)
L:
	if len(keys) > *position && keys[*position] != "" {
		name := keys[*position]
		resource := stns.Attributes{
			name: list[name],
		}

		*position++
		result := s(resource, entry, presult)

		// lack of data
		if result == NSS_STATUS_NOTFOUND {
			goto L
		}

		return result
	}
	return NSS_STATUS_NOTFOUND
}

func (n *Nss) PresetList(list stns.Attributes, position *int) int {
	var attr stns.Attributes
	var err error
	// reset value
	n.PurgeList(list, position)

	attr, err = n.Get()

	if err != nil {
		log.Print(err)
		// When the error refers to the last result.
		// This is supposed to when the server is restarted
		attr = *cache.LastResultList(n.rtype)
		if attr != nil {
			goto C
		}
		return NSS_STATUS_UNAVAIL
	}
	cache.SaveResultList(n.rtype, attr)
C:
	if len(attr) > 0 {
		for k, v := range attr {
			list[k] = v
		}
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_NOTFOUND
}

func (n *Nss) PurgeList(list stns.Attributes, position *int) {
	*position = 0
	for k, _ := range list {
		delete(list, k)
	}
}

func (n *Nss) keys(list stns.Attributes) []string {
	ks := []string{}
	for k, _ := range list {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

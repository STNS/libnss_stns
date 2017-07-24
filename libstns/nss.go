package libstns

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/cache"
	"github.com/STNS/libnss_stns/settings"
)

// http://www.gnu.org/software/libc/manual/html_node/NSS-Modules-Interface.html
const (
	NSS_STATUS_SUCCESS  = 1
	NSS_STATUS_NOTFOUND = 0
	NSS_STATUS_UNAVAIL  = -1
	NSS_STATUS_TRYAGAIN = -2
)

const (
	NSS_LIST_PRESET = 0
	NSS_LIST_PURGE  = 1
)

type Nss struct {
	config  *Config
	rtype   string
	list    stns.Attributes
	listPos *int
}

type NssEntry interface {
	Set(stns.Attributes) int
}

func NewNss(config *Config, rtype string, list stns.Attributes, position *int) *Nss {
	return &Nss{
		config:  config,
		rtype:   rtype,
		list:    list,
		listPos: position,
	}
}

func (n *Nss) Get(column, value string) (stns.Attributes, error) {
	ne := fmt.Errorf("resource notfound %s/%s/%s", n.rtype, column, value)

	req, err := NewRequest(n.config, n.rtype, column, value)
	if err != nil {
		return nil, err
	}

	u, err := cache.Read(req.ApiPath)

	if u != nil || err != nil {
		return u, err
	}

	// default negative cache
	cache.Write(req.ApiPath, nil, ne)
	res, err := req.GetByWrapperCmd()
	if err != nil {
		return nil, err
	}

	if res.Items == nil {
		return nil, ne
	}
	cache.Write(req.ApiPath, res.Items, nil)
	return res.Items, nil
}

func (n *Nss) Set(s NssEntry, column, value string) int {
	id, _ := strconv.Atoi(value)
	minID := cache.ReadMinID(n.rtype)
	maxID := cache.ReadMaxID(n.rtype)

	if n.rtype == "user" {
		if minID+n.config.UIDShift > settings.MIN_LIMIT_ID {
			minID += n.config.UIDShift
		}

		if maxID+n.config.UIDShift > settings.MIN_LIMIT_ID {
			maxID += n.config.UIDShift
		}
	} else if n.rtype == "group" {
		if minID+n.config.GIDShift > settings.MIN_LIMIT_ID {
			minID += n.config.GIDShift
		}

		if maxID+n.config.GIDShift > settings.MIN_LIMIT_ID {
			maxID += n.config.GIDShift
		}
	}
	if column != "id" || (maxID == 0 || minID == 0 || (minID <= id && maxID >= id)) {
		resource, err := n.Get(column, value)
		if err != nil {
			log.Print(err)
			return NSS_STATUS_UNAVAIL
		}

		if len(resource) > 0 {
			return s.Set(resource)
		}
	}
	return NSS_STATUS_NOTFOUND
}

func (n *Nss) SetByList(s NssEntry) int {
	keys := n.keys()
L:
	if len(keys) > *n.listPos && keys[*n.listPos] != "" {
		name := keys[*n.listPos]
		resource := stns.Attributes{
			name: n.list[name],
		}

		*n.listPos++
		result := s.Set(resource)

		// lack of data in list
		// 構造体の項目が不足している場合があるため、配列全て処理を行う
		if result == NSS_STATUS_NOTFOUND {
			goto L
		}
		return result
	}
	return NSS_STATUS_NOTFOUND
}

func (n *Nss) PresetList() int {
	var attr stns.Attributes
	var err error
	// reset value
	n.PurgeList()

	attr, err = n.Get("list", "")

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
			n.list[k] = v
		}
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_NOTFOUND
}

func (n *Nss) PurgeList() {
	*n.listPos = 0
	n.list = stns.Attributes{}
}

func (n *Nss) keys() []string {
	ks := []string{}
	for k, _ := range n.list {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

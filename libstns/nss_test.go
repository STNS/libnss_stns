package libstns

import (
	"testing"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/cache"
)

type Dummy struct {
	t *testing.T
}

func (s Dummy) Set(a stns.Attributes) int {
	for _, u := range a {
		if u.Id != 2000 {
			s.t.Errorf("get error id %d", u.Id)
		}

		if u.GroupId != 3000 {
			s.t.Errorf("get error id %d", u.Id)
		}
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_UNAVAIL
}

func TestGetSuccess(t *testing.T) {
	cache.Flush()
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_01"

	list := stns.Attributes{}
	pos := 0
	n, _ := NewNss(c, "user", list, &pos)

	a, err := n.Get("id", "2000")

	if err != nil {
		t.Errorf("get error %s", err)
	}

	for _, u := range a {
		if u.Id != 2000 {
			t.Errorf("get error id %d", u.Id)
		}
	}

	if cache.ReadMinId("user") != 2000 {
		t.Errorf("get error min_id %d", cache.ReadMinId("user"))
	}

	ca, err := cache.Read("user/id/2000")
	if err != nil {
		t.Errorf("get error %s", err)
	}

	for _, u := range ca {
		if u.Id != 2000 {
			t.Errorf("get error id %d", u.Id)
		}
	}
}

func TestGetNotFound(t *testing.T) {
	cache.Flush()
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_03"

	list := stns.Attributes{}
	pos := 0
	n, _ := NewNss(c, "user", list, &pos)

	a, err := n.Get("id", "2000")
	if a != nil {
		t.Error("get notfound no error 1")
	}
	if err == nil {
		t.Error("get notfound no error 2")
	}
}

func TestSetSuccess(t *testing.T) {
	cache.Flush()
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_01"

	list := stns.Attributes{}
	pos := 0

	d := Dummy{t}

	n, _ := NewNss(c, "user", list, &pos)
	if n.Set(d, "name", "example") != NSS_STATUS_SUCCESS {
		t.Error("set error 1")
	}
}

func TestSetNotFound(t *testing.T) {
	cache.Flush()
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_02"

	list := stns.Attributes{}
	pos := 0

	d := Dummy{t}

	n, _ := NewNss(c, "user", list, &pos)
	r := n.Set(d, "name", "example")
	if r != NSS_STATUS_UNAVAIL {
		t.Errorf("set notfound error response:%d", r)
	}
}

func TestSetByList(t *testing.T) {
	cache.Flush()
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_01"

	list := stns.Attributes{
		"example": &stns.Attribute{
			Id: 2000,
			User: &stns.User{
				GroupId: 3000,
			},
		},
	}
	pos := 0

	d := Dummy{t}

	n, _ := NewNss(c, "user", list, &pos)

	s := n.SetByList(d)
	if s != NSS_STATUS_SUCCESS {
		t.Errorf("set by list error %d", s)
	}

	e := n.SetByList(d)
	if e != NSS_STATUS_NOTFOUND {
		t.Errorf("set by list error %d", e)
	}
}

func TestPresetAndPurgeList(t *testing.T) {
	cache.Flush()
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_01"

	list := stns.Attributes{}
	pos := 0
	n, _ := NewNss(c, "user", list, &pos)
	// preset
	{
		n.PresetList()
		for _, u := range list {
			if u.Id != 2000 {
				t.Errorf("get error id %d", u.Id)
			}
		}

		// cache test
		attr := *cache.LastResultList("user")
		for _, u := range attr {
			if u.Id != 2000 {
				t.Errorf("get error id %d", u.Id)
			}
		}
	}
	// purge
	{
		n.PurgeList()
		if len(list) != 0 {
			t.Errorf("get error list length %d", len(list))
		}
		if pos != 0 {
			t.Errorf("get error list read position %d", pos)
		}

		// rewrite list
		n.PresetList()
		for _, u := range list {
			if u.Id != 2000 {
				t.Errorf("get error id %d", u.Id)
			}
		}
	}
}

package libstns

import (
	"testing"

	"github.com/STNS/libnss_stns/test"
)

// UID,GIDのシフトをテストする
// 通常のレスポンステストはRequestにて包括的にテストしている
func TestUserConvertV3Format(t *testing.T) {
	c := &Config{}
	r, err := convertV3Format([]byte(test.GetV3UserExample()), "/user/name/example", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 2000 {
			t.Error("unmatch id")
		}

		if u.GroupID != 3000 {
			t.Error("unmatch group")
		}
	}

	c = &Config{
		UIDShift: -1000,
		GIDShift: -2000,
	}

	r, err = convertV3Format([]byte(test.GetV3UserExample()), "/user/name/example", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 1000 {
			t.Error("unmatch id")
		}

		if u.GroupID != 1000 {
			t.Error("unmatch group")
		}
	}
}

func TestUsersConvertV3Format(t *testing.T) {
	c := &Config{}
	r, err := convertV3Format([]byte(test.GetV3UsersExample()), "/user/list", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 2000 {
			t.Error("unmatch id")
		}

		if u.GroupID != 3000 {
			t.Error("unmatch group")
		}
	}

	c = &Config{
		UIDShift: -1000,
		GIDShift: -2000,
	}

	r, err = convertV3Format([]byte(test.GetV3UsersExample()), "/user/list", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 1000 {
			t.Error("unmatch id")
		}

		if u.GroupID != 1000 {
			t.Error("unmatch group")
		}
	}
}

func TestGroupConvertV3Format(t *testing.T) {
	c := &Config{}
	r, err := convertV3Format([]byte(test.GetV3GroupExample()), "/group/name/example", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 2000 {
			t.Error("unmatch id")
		}
	}

	c = &Config{
		UIDShift: -1000,
		GIDShift: -1000,
	}

	r, err = convertV3Format([]byte(test.GetV3GroupExample()), "/group/name/example", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 1000 {
			t.Error("unmatch id")
		}
	}
}

func TestGroupsConvertV3Format(t *testing.T) {
	c := &Config{}
	r, err := convertV3Format([]byte(test.GetV3GroupsExample()), "/group/list", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 2000 {
			t.Error("unmatch id")
		}
	}

	c = &Config{
		UIDShift: -1000,
		GIDShift: -1000,
	}

	r, err = convertV3Format([]byte(test.GetV3GroupsExample()), "/group/list", "1000", c)
	if err != nil {
		t.Errorf("not assume an error but got %s", err)
	}

	for n, u := range r.Items {
		if n != "example" {
			t.Error("unmatch name")
		}

		if u.ID != 1000 {
			t.Error("unmatch id")
		}
	}
}

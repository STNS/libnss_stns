package request

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/test"
)

type Response struct {
	path, query, contenttype, body string
}

func TestRequest(t *testing.T) {
	handler := test.GetHandler(t,
		"/user/name/example",
		`{
			"example": {
				"id": 1,
				"group_id": 2,
				"directory": "/home/example",
				"shell": "/bin/sh",
				"keys": [
					"test"
				],
				"password": "password"
			}
		}`,
	)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &config.Config{}
	c.ApiEndPoint = []string{server.URL}
	r, _ := NewRequest(c, "user", "name", "example")

	users, _ := r.Get()
	if 0 == len(users) {
		t.Error("fetch error")
	}

	for n, u := range users {
		if n != "example" {
			t.Error("unmatch name")
		}
		if u.Id != 1 {
			t.Error("unmatch id")
		}
		if u.GroupId != 2 {
			t.Error("unmatch group")
		}
		if u.Directory != "/home/example" {
			t.Error("unmatch direcotry")
		}
		if u.Shell != "/bin/sh" {
			t.Error("unmatch shell")
		}
		if u.Keys[0] != "test" || len(u.Keys) != 1 {
			t.Error("unmatch shell")
		}
		if u.Password != "password" {
			t.Error("unmatch password")
		}
	}

}

func TestBasicAuth(t *testing.T) {
	handler := getBasicAuthHandler(t, "example")
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &config.Config{}
	c.ApiEndPoint = []string{server.URL}
	c.User = "test_user"
	c.Password = "test_pass"
	r, _ := NewRequest(c, "user", "name", "example")

	users, _ := r.Get()
	if 0 == len(users) {
		t.Error("fetch error")
	}

	e := &config.Config{}
	e.ApiEndPoint = []string{server.URL}
	e.User = "error_user"
	e.Password = "error_pass"
	r, _ = NewRequest(e, "user", "name", "example")
	users, _ = r.Get()
	if 0 != len(users) {
		t.Error("fetch error")
	}

}

func TestLockfile(t *testing.T) {
	handler := test.GetHandler(t, "dummy", "dummy")
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &config.Config{}
	c.ApiEndPoint = []string{"example1", "example2"}
	r, _ := NewRequest(c, "dummy", "dummy", "dummy")

	r.Get()
	lock1 := "/tmp/libnss_stns." + GetMD5Hash("example1")
	lock2 := "/tmp/libnss_stns." + GetMD5Hash("example2")

	_, err := os.Stat(lock1)
	if err != nil {
		t.Error("not exist lock file 1")
	}

	_, err = os.Stat(lock2)
	if err != nil {
		t.Error("not exist lock file 2")
	}
}

func getBasicAuthHandler(t *testing.T, name string) http.HandlerFunc {
	response := &Response{
		path:        "/user/name/" + name,
		contenttype: "application/json",
		body: `{
"example": {
	"id": 1,
	"group_id": 2,
	"directory": "/home/example",
	"shell": "/bin/sh",
	"keys": [
		"test"
	]
}
}`,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request.
		if g, w := r.URL.Path, response.path; g != w {
			t.Errorf("request got path %s, want %s", g, w)
		}
		if authName, authPass, authStatus := r.BasicAuth(); authStatus {
			if authName != "test_user" || authPass != "test_pass" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
		}

		// Send response.
		w.Header().Set("Content-Type", response.contenttype)
		io.WriteString(w, response.body)
	}
	return handler
}

package request

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func TestRequest(t *testing.T) {
	handler := getHandler(t, "example")
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ConfigFileName = ""
	r, _ := NewRequest("user", "name", "example")
	r.Config.ApiEndPoint = []string{server.URL}
	r.Config.User = "test_user"
	r.Config.Password = "test_pass"

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
	}
}

func TestErrorBasicAuth(t *testing.T) {
	handler := getHandler(t, "example")
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ConfigFileName = ""
	r, _ := NewRequest("user", "name", "example")
	r.Config.ApiEndPoint = []string{server.URL}
	r.Config.User = "error_user"
	r.Config.Password = "error_pass"
	Cache = map[string]*CacheObject{}
	users, _ := r.Get()
	if 0 != len(users) {
		t.Error("fetch error")
	}

}

func getHandler(t *testing.T, name string) http.HandlerFunc {
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

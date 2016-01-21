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
	r.Config.ApiEndPoint = server.URL
	users, _ := r.Get()
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

		// Send response.
		w.Header().Set("Content-Type", response.contenttype)
		io.WriteString(w, response.body)
	}
	return handler
}

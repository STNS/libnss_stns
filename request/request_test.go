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
	attr, _ := Send(server.URL + "/user/name/example")

	if attr.Id != 1 {
		t.Error("unmatch id")
	}

	if attr.Name != "example" {
		t.Error("unmatch name")
	}
	if attr.GroupId != 2 {
		t.Error("unmatch group")
	}
	if attr.Directory != "/home/example" {
		t.Error("unmatch direcotry")
	}
	if attr.Shell != "/bin/sh" {
		t.Error("unmatch shell")
	}
	if attr.Keys[0] != "test" || len(attr.Keys) != 1 {
		t.Error("unmatch shell")
	}
}

func getHandler(t *testing.T, name string) http.HandlerFunc {
	response := &Response{
		path:        "/user/name/" + name,
		contenttype: "application/json",
		body: `{
"meta": {
"code": 200
},
"id": 1,
"name": "example",
"group_id": 2,
"directory": "/home/example",
"shell": "/bin/sh",
"keys": [
	"test"
]
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

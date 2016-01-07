package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pyama86/libnss_stns/internal"
)

type Response struct {
	path, query, contenttype, body string
}

func TestFetchKey(t *testing.T) {
	var config libnss_stns.Config
	okhandler := getHandler(t, "example", `"test key1","test key2"`)
	okserver := httptest.NewServer(http.HandlerFunc(okhandler))
	config.ApiEndPoint = okserver.URL
	defer okserver.Close()
	if "test key1\ntest key2" != FetchKey("example", &config) {
		t.Error("unmatch keys")
	}

	nghandler := getHandler(t, "notfound", "")
	ngserver := httptest.NewServer(http.HandlerFunc(nghandler))
	config.ApiEndPoint = ngserver.URL
	defer ngserver.Close()

	if "" != FetchKey("notfound", &config) {
		t.Error("unmatch keys")
	}

}

func getHandler(t *testing.T, name string, keys string) http.HandlerFunc {
	response := &Response{
		path:        "/user/name/" + name,
		contenttype: "application/json",
		body: fmt.Sprintf(`{
"meta": {
"code": 200
},
"id": 1,
"name": "example",
"group_id": 2,
"directory": "/home/example",
"shell": "/bin/sh",
"keys": [
	%s
]
}`, keys),
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

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func TestFetchKey(t *testing.T) {
	//t.Errorf("Root(confing from file) should be /hoge/fuga but: %v", "a")

	response := &Response{
		path:        "/user/name/example",
		contenttype: "application/json",
		body: `{
"meta": {
"code": 400
},
"id": 1,
"name": "example",
"group_id": 2,
"directory": "/home/example",
"shell": "/bin/sh",
"keys": [
	"test key"
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

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	if "test key" != FetchKey("example", server.URL) {
		t.Error("unmatch keys")
	}

}

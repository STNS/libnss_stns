package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/STNS/libnss_stns/config"
)

type Response struct {
	path, query, contenttype, body string
}

func TestFetchKey(t *testing.T) {
	okhandler := GetHandler(t)
	okserver := httptest.NewServer(http.HandlerFunc(okhandler))
	c := &config.Config{ApiEndPoint: []string{okserver.URL}}
	defer okserver.Close()
	if "test key1\ntest key2" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

	nghandler := NgGetHandler(t)
	ngserver := httptest.NewServer(http.HandlerFunc(nghandler))
	c = &config.Config{ApiEndPoint: []string{ngserver.URL}}
	defer ngserver.Close()

	if "" != Fetch(c, "notfound") {
		t.Error("unmatch keys")
	}

	// fail over
	c = &config.Config{ApiEndPoint: []string{"", okserver.URL}}
	defer ngserver.Close()

	if "test key1\ntest key2" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

}

func GetHandler(t *testing.T) http.HandlerFunc {
	response := &Response{
		path:        "/user/name/example",
		contenttype: "application/json",
		body: `{
"example": {
	"id": 1,
	"group_id": 2,
	"directory": "/home/example",
	"shell": "/bin/sh",
	"keys": [
		"test key1",
		"test key2"
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

func NgGetHandler(t *testing.T) http.HandlerFunc {
	response := &Response{
		path:        "/user/name/notfound",
		contenttype: "application/json",
		body: `{
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

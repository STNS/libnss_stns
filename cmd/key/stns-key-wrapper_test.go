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
	successHandler := GetHandler(t,
		"/user/name/example",
		`{
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
	)
	successServer := httptest.NewServer(http.HandlerFunc(successHandler))
	c := &config.Config{ApiEndPoint: []string{successServer.URL}}
	defer successServer.Close()

	if "test key1\ntest key2" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

	notfoundHandler := GetHandler(t,
		"/user/name/notfound",
		`{
		}`,
	)
	notfoundServer := httptest.NewServer(http.HandlerFunc(notfoundHandler))
	c = &config.Config{ApiEndPoint: []string{notfoundServer.URL}}
	defer notfoundServer.Close()

	if "" != Fetch(c, "notfound") {
		t.Error("unmatch keys")
	}

	// fail over
	c = &config.Config{ApiEndPoint: []string{"", successServer.URL}}
	defer notfoundServer.Close()

	if "test key1\ntest key2" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

}

func GetHandler(t *testing.T, path string, responseBody string) http.HandlerFunc {
	response := &Response{
		path:        path,
		contenttype: "application/json",
		body:        responseBody,
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

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/STNS/libnss_stns/config"
)

type Response struct {
	path, query, contenttype, body string
}

func useTestBins(t *testing.T) func() {
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "./test-fixtures/bin:/bin:/usr/bin")
	return func() { os.Setenv("PATH", origPath) }
}

func TestFetchKey(t *testing.T) {
	// normal
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
	defer successServer.Close()

	c := &config.Config{ApiEndPoint: []string{successServer.URL}}

	if "test key1\ntest key2\n" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

	// user notfound
	notfoundHandler := GetHandler(t,
		"/user/name/notfound",
		`{
		}`,
	)
	notfoundServer := httptest.NewServer(http.HandlerFunc(notfoundHandler))
	defer notfoundServer.Close()

	c = &config.Config{ApiEndPoint: []string{notfoundServer.URL}}

	if "" != Fetch(c, "notfound") {
		t.Error("unmatch keys")
	}

	// fail over
	c = &config.Config{ApiEndPoint: []string{"", successServer.URL}}

	if "test key1\ntest key2\n" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

	// chain wrapper
	{
		defer useTestBins(t)()

		c = &config.Config{ApiEndPoint: []string{successServer.URL}, ChainSshWrapper: "get-external-keys"}
		if "test key1\ntest key2\nexternal key1\nexternal key2\n" != Fetch(c, "example") {
			t.Errorf("unmatch keys: '%#v'", Fetch(c, "example"))
		}
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

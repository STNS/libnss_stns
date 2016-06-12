package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/STNS/libnss_stns/libstns"
	"github.com/STNS/libnss_stns/test"
)

func TestFetch(t *testing.T) {
	// normal
	successHandler := test.GetHandler(t,
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
		200,
	)
	successServer := httptest.NewServer(http.HandlerFunc(successHandler))
	defer successServer.Close()

	r := regexp.MustCompile(`example`)
	c := &libstns.Config{ApiEndPoint: []string{successServer.URL}}
	out, _ := Fetch(c, "/user/name/example")
	if !r.MatchString(out) {
		t.Error("unmatch response")
	}

	// user notfound
	r = regexp.MustCompile(`{\s+}`)
	notfoundHandler := test.GetHandler(t,
		"/user/name/notfound",
		`{
		}`,
		404,
	)
	notfoundServer := httptest.NewServer(http.HandlerFunc(notfoundHandler))
	defer notfoundServer.Close()

	c = &libstns.Config{ApiEndPoint: []string{notfoundServer.URL}}

	out, _ = Fetch(c, "/user/name/notfound")
	if r.MatchString(out) {
		t.Error("unmatch keys")
	}
}

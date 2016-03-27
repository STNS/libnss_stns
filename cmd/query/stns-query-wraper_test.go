package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/STNS/lib-stns/config"
	"github.com/STNS/lib-stns/test"
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
	)
	successServer := httptest.NewServer(http.HandlerFunc(successHandler))
	defer successServer.Close()

	r := regexp.MustCompile(`example`)
	c := &config.Config{ApiEndPoint: []string{successServer.URL}}

	if !r.MatchString(Fetch(c, "/user/name/example")) {
		t.Error("unmatch response")
	}

	// user notfound
	r = regexp.MustCompile(`{\s+}`)
	notfoundHandler := test.GetHandler(t,
		"/user/name/notfound",
		`{
		}`,
	)
	notfoundServer := httptest.NewServer(http.HandlerFunc(notfoundHandler))
	defer notfoundServer.Close()

	c = &config.Config{ApiEndPoint: []string{notfoundServer.URL}}
	fmt.Println(Fetch(c, "/user/name/notfound"))
	if !r.MatchString(Fetch(c, "/user/name/notfound")) {
		t.Error("unmatch keys")
	}
}

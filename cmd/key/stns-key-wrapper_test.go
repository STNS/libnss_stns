package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/STNS/libnss_stns/libstns"
	"github.com/STNS/libnss_stns/test"
)

func useTestBins(t *testing.T) func() {
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "./test-fixtures/bin:/bin:/usr/bin")
	return func() { os.Setenv("PATH", origPath) }
}

func TestFetchKey(t *testing.T) {
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

	c := &libstns.Config{ApiEndPoint: []string{successServer.URL}}

	result := Fetch(c, "example")
	if "test key1\ntest key2\n" != result {
		t.Errorf("unmatch key value %s", result)
	}

	// no key
	nokeyHandler := test.GetHandler(t,
		"/user/name/example",
		`{
			"example": {
				"id": 1,
				"group_id": 2,
				"directory": "/home/example",
				"shell": "/bin/sh",
				"keys": []
			}
		}`,
		404,
	)
	nokeyServer := httptest.NewServer(http.HandlerFunc(nokeyHandler))
	defer nokeyServer.Close()

	c = &libstns.Config{ApiEndPoint: []string{nokeyServer.URL}}
	if "" != Fetch(c, "example") {
		t.Error("unmatch keys by nokey")
	}

	// user notfound
	notfoundHandler := test.GetHandler(t,
		"/user/name/notfound",
		`{
		}`,
		404,
	)

	notfoundServer := httptest.NewServer(http.HandlerFunc(notfoundHandler))
	defer notfoundServer.Close()

	c = &libstns.Config{ApiEndPoint: []string{notfoundServer.URL}}

	if "" != Fetch(c, "notfound") {
		t.Error("unmatch keys")
	}

	// fail over
	c = &libstns.Config{ApiEndPoint: []string{"", successServer.URL}}

	if "test key1\ntest key2\n" != Fetch(c, "example") {
		t.Error("unmatch keys")
	}

	// chain wrapper
	{
		defer useTestBins(t)()

		c = &libstns.Config{ApiEndPoint: []string{successServer.URL}, ChainSshWrapper: "get-external-keys"}
		if "test key1\ntest key2\nexternal key1\nexternal key2\n" != Fetch(c, "example") {
			t.Errorf("unmatch keys: '%#v'", Fetch(c, "example"))
		}
	}
}

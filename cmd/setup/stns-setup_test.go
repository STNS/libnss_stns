package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/STNS/libnss_stns/libstns"
	"github.com/STNS/libnss_stns/test"
)

func TestRun(t *testing.T) {
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
				],
				"setup_commands": [
					"./test-fixtures/bin/dummy.sh"
				]
			}
		}`,
		200,
	)
	successServer := httptest.NewServer(http.HandlerFunc(successHandler))
	defer successServer.Close()

	c := &libstns.Config{ApiEndPoint: []string{successServer.URL}}
	err := Run(c, "example")
	if err != nil {
		t.Errorf("did not assume an error: %s", err)
	}
}

package test

import (
	"io"
	"net/http"
	"testing"
)

type Response struct {
	path, query, contenttype, body string
}

func GetHandler(t *testing.T, path string, responseBody string, responseCode int) http.HandlerFunc {
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
		w.WriteHeader(responseCode)
		io.WriteString(w, response.body)
	}
	return handler
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

func GetV1Example() string {
	return `{
		"example": {
			"id": 2000,
		"group_id": 3000,
			"directory": "/home/example",
			"shell": "/bin/sh",
			"keys": [
				"test"
			],
			"password": "password"
		}
	}`
}
func GetV2Example() string {
	return `{
		"metadata": {
			"api_version": 2.0,
			"salt_enable": false,
			"stretching_number": 0,
			"result": "success",
			"min_id": 1000
		},
		"items": {
			"example": {
				"id": 2000,
				"group_id": 3000,
				"directory": "/home/example",
				"shell": "/bin/sh",
				"keys": [
					"test"
				],
				"password": "password"
			}
		}
	}`

}

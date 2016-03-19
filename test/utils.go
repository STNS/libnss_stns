package test

import (
	"io"
	"net/http"
	"testing"
)

type Response struct {
	path, query, contenttype, body string
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

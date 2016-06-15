package libstns

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/test"
)

type Response struct {
	path, query, contenttype, body string
}

func TestRequestV1ServerV1(t *testing.T) {
	handler := test.GetHandler(t,
		"/user/name/example",
		test.GetV1Example(),
		200,
	)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, 1.0)

}

func TestRequestV2ServerV2(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", test.GetV2Example(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, 2.0)
}

func TestRequestV2NotFound(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", `{
	"metadata": {
		"api_version": 2.0,
		"salt_enable": false,
		"stretching_number": 0,
		"result": "success"
	},
	"items": null
	}`,
		404,
	)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	var res stns.ResponseFormat
	raw, err := r.GetRawData()
	json.Unmarshal(raw, &res)
	if err != nil {
		t.Errorf("fetch error %s", err)
	}
}

func TestFailOver(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", test.GetV2Example(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{"http://localhost:1000", server.URL + "/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, 2.0)
}

func TestRefused(t *testing.T) {
	c := &Config{}
	c.ApiEndPoint = []string{"http://localhost:1000"}

	r, _ := NewRequest(c, "user", "name", "example")
	_, err := r.GetRawData()
	if err == nil {
		t.Error("errot test refused")
	}
}

func checkAttribute(t *testing.T, res stns.ResponseFormat, apiVersion float64) {
	// metadata
	if res.MetaData.ApiVersion != apiVersion {
		t.Error("unmatch api version")
	}

	if res.MetaData.Salt {
		t.Error("unmatch salt")
	}

	if res.MetaData.Stretching != 0 {
		t.Error("unmatch stretching")
	}

	if res.MetaData.Result != "success" {
		t.Error("unmatch result")
	}

	if res.MetaData.ApiVersion == 2.0 {
		if res.MetaData.MinId != 2000 {
			t.Errorf("unmatch min id %d", res.MetaData.MinId)
		}
	}

	for n, u := range *res.Items {
		if n != "example" {
			t.Error("unmatch name")
		}
		if u.Id != 2000 {
			t.Error("unmatch id")
		}
		if u.GroupId != 3000 {
			t.Error("unmatch group")
		}
		if u.Directory != "/home/example" {
			t.Error("unmatch direcotry")
		}
		if u.Shell != "/bin/sh" {
			t.Error("unmatch shell")
		}
		if u.Keys[0] != "test" || len(u.Keys) != 1 {
			t.Error("unmatch keys")
		}
		if u.Password != "password" {
			t.Error("unmatch password")
		}
	}
}

func checkResponse(t *testing.T, r *Request, apiVersion float64) {
	var res stns.ResponseFormat
	raw, err := r.GetRawData()
	json.Unmarshal(raw, &res)
	if err != nil || res.Items == nil || 0 == len(*res.Items) {
		t.Errorf("fetch error %s", err)
	}
	if err == nil {
		checkAttribute(t, res, apiVersion)
	}
}

func TestBasicAuth(t *testing.T) {
	handler := getBasicAuthHandler(t, "example")
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL}
	c.User = "test_user"
	c.Password = "test_pass"
	r, _ := NewRequest(c, "user", "name", "example")

	users, _ := r.GetRawData()
	if 0 == len(users) {
		t.Error("fetch error")
	}

	e := &Config{}
	e.ApiEndPoint = []string{server.URL}
	e.User = "error_user"
	e.Password = "error_pass"
	r, _ = NewRequest(e, "user", "name", "example")
	users, _ = r.GetRawData()
	if 0 != len(users) {
		t.Error("fetch error")
	}

}

func getBasicAuthHandler(t *testing.T, name string) http.HandlerFunc {
	response := &Response{
		path:        "/user/name/" + name,
		contenttype: "application/json",
		body: `{
"example": {
	"id": 1,
	"group_id": 2,
	"directory": "/home/example",
	"shell": "/bin/sh",
	"keys": [
		"test"
	]
}
}`,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check request.
		if g, w := r.URL.Path, response.path; g != w {
			t.Errorf("request got path %s, want %s", g, w)
		}
		if authName, authPass, authStatus := r.BasicAuth(); authStatus {
			if authName != "test_user" || authPass != "test_pass" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}
		}

		// Send response.
		w.Header().Set("Content-Type", response.contenttype)
		io.WriteString(w, response.body)
	}
	return handler
}

func TestGetByWrapperCmd(t *testing.T) {
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_01"
	r, _ := NewRequest(c, "user", "name", "example")
	res, _ := r.GetByWrapperCmd()
	checkAttribute(t, res, 2.0)
}

func TestGetByWrapperCmd404(t *testing.T) {
	c := &Config{}
	c.ApiEndPoint = []string{"exmple"}
	c.WrapperCommand = "./fixtures/bin/command_response_02"
	r, _ := NewRequest(c, "user", "name", "example")
	_, err := r.GetByWrapperCmd()
	if err != nil {
		t.Errorf("fetch error %s", err)
	}
}

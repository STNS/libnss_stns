package libstns

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/STNS/libnss_stns/cache"
	"github.com/STNS/libnss_stns/test"
)

type Response struct {
	path, query, contenttype, body string
}

func TestRequestTimeOut(t *testing.T) {
	c := &Config{}
	c.ApiEndPoint = []string{"http://10.1.1.1/v2"}
	c.RequestTimeOut = 1
	r, _ := NewRequest(c, "user", "name", "example")
	_, err := r.GetRawData()

	if err == nil {
		t.Error("fetch timeout error")
	}
}

func TestRequestProxyByEnv(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", test.GetV2Example(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	os.Setenv("HTTP_PROXY", server.URL+"/v2")

	defer server.Close()
	defer os.Unsetenv("HTTP_PROXY")

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, checkUserAttribute)
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
	checkResponse(t, r, checkUserAttribute)

}

func TestRequestV2ServerV2(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", test.GetV2Example(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, checkUserAttribute)
}

func TestRequestV3ServerV3User(t *testing.T) {
	handler := test.GetHandler(t, "/v3/user/name/example", test.GetV3UserExample(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v3"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, checkUserAttribute)
}

func TestRequestV3ServerV3Users(t *testing.T) {
	handler := test.GetHandler(t, "/v3/user/list", test.GetV3UsersExample(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v3"}

	r, _ := NewRequest(c, "user", "list", "")
	checkResponse(t, r, checkUserAttribute)
}

func TestRequestV3ServerV3Group(t *testing.T) {
	handler := test.GetHandler(t, "/v3/group/name/example", test.GetV3GroupExample(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v3"}

	r, _ := NewRequest(c, "group", "name", "example")
	checkResponse(t, r, checkGroupAttribute)
}

func TestRequestV3ServerV3Groups(t *testing.T) {
	handler := test.GetHandler(t, "/v3/group/list", test.GetV3GroupsExample(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v3"}

	r, _ := NewRequest(c, "group", "list", "")
	checkResponse(t, r, checkGroupAttribute)
}

func TestRequestV3ServerV3Sudo(t *testing.T) {
	handler := test.GetHandler(t, "/v3/sudo/name/example", test.GetV3UserExample(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v3"}

	r, _ := NewRequest(c, "sudo", "name", "example")
	checkResponse(t, r, checkSudoAttribute)
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
	var res ResponseFormat
	raw, err := r.GetRawData()
	json.Unmarshal(raw, &res)
	if err == nil {
		t.Errorf("fetch error expect not found")
	}
}

func TestRequestV3NotFound(t *testing.T) {
	handler := test.GetHandler(t, "/v3/user/name/example", `{}`, 404)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL + "/v3"}

	r, _ := NewRequest(c, "user", "name", "example")
	var res ResponseFormat
	raw, err := r.GetRawData()
	json.Unmarshal(raw, &res)
	if err == nil {
		t.Errorf("fetch error expect not found")
	}
}

func TestFailOver(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", test.GetV2Example(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{"http://localhost:1000", server.URL + "/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, checkUserAttribute)
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
	checkUserAttribute(t, res)
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

func TestRequestProxyByConfig(t *testing.T) {
	handler := test.GetHandler(t, "/v2/user/name/example", test.GetV2Example(), 200)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}

	c.HttpProxy = server.URL + "/v2"
	c.ApiEndPoint = []string{"http://unservice_config/v2"}

	r, _ := NewRequest(c, "user", "name", "example")
	checkResponse(t, r, checkUserAttribute)
}

func checkUserAttribute(t *testing.T, res *ResponseFormat) {

	for n, u := range res.Items {
		if n != "example" {
			t.Error("unmatch name")
		}
		if u.ID != 2000 {
			t.Error("unmatch id")
		}
		if u.GroupID != 3000 {
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

func checkGroupAttribute(t *testing.T, res *ResponseFormat) {
	for n, g := range res.Items {
		if n != "example" {
			t.Error("unmatch name")
		}
		if g.ID != 2000 {
			t.Error("unmatch id")
		}
		if g.Users[0] != "test" || len(g.Users) != 1 {
			t.Error("unmatch users")
		}
	}
}

func checkSudoAttribute(t *testing.T, res *ResponseFormat) {
	for n, u := range res.Items {
		if n != "example" {
			t.Error("unmatch name")
		}
		if u.Password != "password" {
			t.Error("unmatch password")
		}
	}
}

func checkResponse(t *testing.T, r *Request, checkAttribute func(*testing.T, *ResponseFormat)) {
	cache.Flush()
	var res ResponseFormat

	raw, err := r.GetRawData()
	if err != nil {
		t.Errorf("fetch error %s", err)
	}

	err = json.Unmarshal(raw, &res)

	if err != nil {
		t.Errorf("fetch error %s", err)
	}

	if res.Items == nil || 0 == len(res.Items) {
		t.Error("fetch error response is nil")
	} else {
		checkAttribute(t, &res)
	}
}

func TestRequestHeader(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if !strings.Contains(apiKey, "test") {
			t.Errorf("unmatch header  error %s", apiKey)
		}
	})

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &Config{}
	c.ApiEndPoint = []string{server.URL}
	c.RequestHeader = map[string]string{"x-api-key": "test"}

	r, _ := NewRequest(c, "user", "name", "example")
	r.GetRawData()
}

func TestTlsAuth(t *testing.T) {
	c := &Config{}
	c.TlsCa = "./fixtures/keys/test.pem"
	c.SslVerify = false
	r, _ := NewRequest(c, "user", "name", "example")

	if r.TlsConfig().InsecureSkipVerify == false {
		t.Error("tls auth error 1")
	}

	if r.TlsConfig().RootCAs != nil {
		t.Error("tls auth error 2")
	}

	c.TlsCert = "./fixtures/keys/test.crt"
	c.TlsKey = "./fixtures/keys/test.key"
	r, _ = NewRequest(c, "user", "name", "example")

	if r.TlsConfig().RootCAs == nil {
		t.Error("tls auth error 3")
	}

	c.TlsCert = "./fixtures/keys/not_exists.crt"
	c.TlsKey = "./fixtures/keys/test.key"
	r, _ = NewRequest(c, "user", "name", "example")

	if r.TlsConfig().RootCAs != nil {
		t.Error("tls auth error 4")
	}
}

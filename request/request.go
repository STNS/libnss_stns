package request

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/lib-stns/config"
	"github.com/STNS/lib-stns/logger"
)

type Request struct {
	ApiPath string
	Config  *config.Config
}

func NewRequest(config *config.Config, paths ...string) (*Request, error) {
	logger.Setlog()
	r := Request{}

	r.Config = config
	r.ApiPath = strings.Join(paths, "/")

	return &r, nil
}

func (r *Request) GetRaw() ([]byte, error) {
	var lastError error
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(r.Config.ApiEndPoint))

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: !r.Config.SslVerify}

	for _, v := range perm {
		endPoint := r.Config.ApiEndPoint[v]
		url := strings.TrimRight(endPoint, "/") + "/" + strings.TrimLeft(path.Clean(r.ApiPath), "/")
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			lastError = err
			continue
		}

		if r.Config.User != "" && r.Config.Password != "" {
			req.SetBasicAuth(r.Config.User, r.Config.Password)
		}

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			lastError = err
			continue
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			lastError = err
			continue
		}

		if res.StatusCode == http.StatusOK {
			return body, nil
		}
	}
	return nil, lastError
}

func (r *Request) Get() (attribute.AllAttribute, error) {
	var attr attribute.AllAttribute

	body, err := r.GetRaw()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &attr)

	if err != nil {
		return nil, err
	}

	return attr, nil
}

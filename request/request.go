package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/libnss_stns/config"
)

var ConfigFileName = "/etc/stns/libnss_stns.conf"
var Loaded *config.Config

type Request struct {
	apiPath string
	Config  *config.Config
}

type HttpResponse struct {
	response *http.Response
	err      error
}

func NewRequest(resource string, column string, value string) (*Request, error) {
	r := Request{}
	if err := r.Init(); err != nil {
		return nil, err
	}
	urls := []string{resource, column}

	if value != "" {
		urls = append(urls, value)
	}

	r.apiPath = strings.Join(urls, "/")
	return &r, nil
}

func (r *Request) asyncHttpGets() ([]byte, error) {
	responses := []*HttpResponse{}
	ch := make(chan *HttpResponse, len(r.Config.ApiEndPoint))
	for _, endPoint := range r.Config.ApiEndPoint {
		go func(endPoint string) {
			url := endPoint + "/" + r.apiPath
			res, err := http.Get(url)

			if err != nil {
				ch <- &HttpResponse{nil, err}
			} else {
				ch <- &HttpResponse{res, err}
			}
		}(endPoint)
	}

	for {
		select {
		case c := <-ch:
			responses = append(responses, c)
			if c.response != nil {
				if c.response.StatusCode != http.StatusOK {
					continue
				}
				body, err := ioutil.ReadAll(c.response.Body)
				if err != nil {
					continue
				}
				c.response.Body.Close()
				return body, nil
			} else if c.response == nil && len(responses) == len(r.Config.ApiEndPoint) {
				return nil, c.err
			}
		}
	}
}
func (r *Request) Get() (attribute.UserGroups, error) {
	body, err := r.asyncHttpGets()
	if err != nil {
		return nil, err
	}

	var attr attribute.UserGroups
	err = json.Unmarshal(body, &attr)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

func (r *Request) Init() error {

	if reflect.ValueOf(Loaded).IsNil() {
		logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, os.Args[0])
		if err != nil {
			// syslog not found
			fmt.Print(err)
		} else {
			log.SetOutput(logger)
		}
		if ConfigFileName != "" {
			config, err := config.Load(ConfigFileName)
			if err != nil {
				log.Print(err)
				return err
			}
			Loaded = config
		} else {
			Loaded = &config.Config{}
		}
	}
	r.Config = Loaded
	return nil
}

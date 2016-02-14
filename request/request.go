package request

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/libnss_stns/config"
)

var ConfigFileName = "/etc/stns/libnss_stns.conf"
var Loaded *config.Config
var Cache map[string]*CacheObject

var Pid int

type CacheObject struct {
	userGroup *attribute.UserGroups
	err       error
}

type Request struct {
	ApiPath string
	Config  *config.Config
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

	r.ApiPath = strings.Join(urls, "/")

	if Pid != 0 && Pid != os.Getpid() {
		return nil, errors.New("unsupported fork process")
	}

	Pid = os.Getpid()
	return &r, nil
}

func (r *Request) Get() (attribute.UserGroups, error) {
	var lastError error
	var attr attribute.UserGroups

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(r.Config.ApiEndPoint))

	c, exist := Cache[r.ApiPath]
	if exist {
		if c.err != nil {
			return nil, c.err
		} else {
			return *c.userGroup, c.err
		}
	}

	// default negative cache
	Cache[r.ApiPath] = &CacheObject{nil, errors.New(r.ApiPath + " is not fond")}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: !r.Config.SslVerify}

	for _, v := range perm {
		endPoint := r.Config.ApiEndPoint[v]
		url := endPoint + "/" + r.ApiPath
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
			err = json.Unmarshal(body, &attr)
			if err != nil {
				lastError = err
				continue
			}
			Cache[r.ApiPath] = &CacheObject{&attr, nil}
			return attr, nil
		}
	}
	return nil, lastError
}

func (r *Request) Init() error {
	if len(Cache) == 0 {
		Cache = map[string]*CacheObject{}
	}
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

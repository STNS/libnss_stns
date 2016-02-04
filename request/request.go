package request

import (
	"encoding/json"
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

type Request struct {
	apiPath string
	Config  *config.Config
}

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(s))
	return s[i]
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

func (r *Request) Get() (attribute.UserGroups, error) {
	var resultErr error
	client := &http.Client{}

	for _, endPoint := range r.Config.ApiEndPoint {
		req, err := http.NewRequest("GET", endPoint+"/"+r.apiPath, nil)
		if err != nil {
			resultErr = err
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			resultErr = err
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			resultErr = err
			continue
		}

		var attr attribute.UserGroups
		err = json.Unmarshal(body, &attr)
		if err != nil {
			resultErr = err
			continue
		}
		return attr, nil
	}
	return nil, resultErr
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

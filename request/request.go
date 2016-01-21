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

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/config"
)

var ConfigFileName = "/etc/stns/libnss_stns.conf"
var Loaded *config.Config

type Request struct {
	Url    string
	Config *config.Config
}

func NewRequest(resource string, column string, value string) (*Request, error) {
	r := Request{}
	if err := r.Init(); err != nil {
		return nil, err
	}

	urls := []string{r.Config.ApiEndPoint, resource, column}

	if value != "" {
		urls = append(urls, value)
	}

	r.Url = strings.Join(urls, "/")
	return &r, nil
}

func (r *Request) Get() (attribute.UserGroups, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, err := ioutil.ReadAll(res.Body)
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

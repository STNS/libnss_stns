package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"reflect"
	"strings"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/config"
)

const configFile = "/etc/stns/libnss_stns.conf"

var Loaded *config.Config

type Request struct {
	Config *config.Config
	PgName string
}

func (r *Request) Get(resources []string) (attribute.UserGroups, error) {
	if err := r.Init(); err != nil {
		return nil, err
	}

	s := []string{r.Config.ApiEndPoint}
	append(s, resources)
	return Send(strings.Join(s, "/"))
}

func Send(url string) (attribute.UserGroups, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
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
		logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, r.PgName)
		if err != nil {
			// syslog not found
			fmt.Print(err)
		} else {
			log.SetOutput(logger)
		}

		config, err := config.Load(configFile)
		if err != nil {
			log.Print(err)
			return err
		}
		Loaded = config
	}
	r.Config = Loaded
	return nil
}

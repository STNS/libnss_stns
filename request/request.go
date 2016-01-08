package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/init"
)

func Get(resource string, column string, value string) (map[string]*attribute.All, error) {
	config, err := libnss_stns.Init()
	if err != nil {
		return nil, err
	}
	s := []string{config.ApiEndPoint, resource, column, value}

	list, err := Send(strings.Join(s, "/"))

	if err != nil {
		return nil, err
	}
	return list, err
}

func GetList(resource string) (map[string]*attribute.All, error) {
	config, err := libnss_stns.Init()
	if err != nil {
		return nil, err
	}
	s := []string{config.ApiEndPoint, resource, "list"}

	list, err := Send(strings.Join(s, "/"))

	if err != nil {
		return nil, err
	}
	return list, err
}

func Send(url string) (map[string]*attribute.All, error) {
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

	var attr map[string]*attribute.All
	err = json.Unmarshal(body, &attr)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

package libnss_stns

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Request(url string) (*Attr, error) {
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

	var attr Attr
	err = json.Unmarshal(body, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

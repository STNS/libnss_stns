package libnss_stns

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
)

func Request(url string) *Attr {
	logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "libnss_stns")
	log.SetOutput(logger)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	var attr Attr
	err = json.Unmarshal(body, &attr)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &attr
}

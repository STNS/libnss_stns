package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

func main() {
	libstns.Setlog()
	config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
	if err == nil {
		flag.Parse()
		if raw := Fetch(config, flag.Arg(0)); raw != "" {
			fmt.Println(raw)
		}
	}
}

func Fetch(config *libstns.Config, name string) string {
	var res stns.ResponseFormat
	var userKeys string

	r, err := libstns.NewRequest(config, "user", "name", name)
	if err != nil {
		log.Println(err)
	}

	body, err := r.GetRawData()
	if err != nil {
		log.Println(err)
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &res)

		if err != nil {
			log.Println(err)
		}
	}

	if res.Items != nil {
		for _, u := range res.Items {
			if len(u.Keys) > 0 {
				userKeys += strings.Join(u.Keys, "\n") + "\n"
			}
		}
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
			}
		}()
	}

	rex := regexp.MustCompile(`stns-key-wrapper$`)
	if r.Config.ChainSshWrapper != "" && !rex.MatchString(r.Config.ChainSshWrapper) {
		out, err := exec.Command(r.Config.ChainSshWrapper, name).Output()
		if err != nil {
			log.Print(err)
		} else {
			if "" != string(out) {
				userKeys += string(out)
			}
		}

	}
	return userKeys
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/logger"
	"github.com/STNS/libnss_stns/request"
)

func main() {
	logger.Setlog()
	config, err := config.Load("/etc/stns/libnss_stns.conf")
	if err == nil {
		flag.Parse()
		if raw := Fetch(config, flag.Arg(0)); raw != "" {
			fmt.Println(raw)
		}
	}
}

func Fetch(config *config.Config, name string) string {
	userKeys := ""
	r, err := request.NewRequest(config, "user", "name", name)
	if err != nil {
		log.Println(err)
	}

	users, err := r.GetAttributes()
	if err != nil {
		log.Println(err)
	}

	if users != nil {
		for _, u := range users {
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

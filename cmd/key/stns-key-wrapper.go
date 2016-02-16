package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/STNS/libnss_stns/cmd"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/request"
)

func main() {
	config, err := cmd.LoadConfig()
	if err == nil {
		if raw := Fetch(config, flag.Arg(0)); raw != "" {
			fmt.Println(raw)
		}
	}
}

func Fetch(config *config.Config, name string) string {
	userKeys := ""
	r, err := request.NewRequest(config, "user", "name", name)
	if err != nil {
		log.Print(err)
		return ""
	}

	users, err := r.Get()
	if err != nil {
		log.Print(err)
	}

	if users != nil {
		for _, u := range users {
			userKeys += strings.Join(u.Keys, "\n")
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
			return ""
		}
		userKeys += string(out)
	}
	return userKeys
}

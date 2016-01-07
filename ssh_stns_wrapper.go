package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/pyama86/libnss_stns/internal"
)

const configFile = "/etc/stns/libnss_stns.conf"

func main() {
	flag.Parse()
	if err := libnss_stns.InitLogger("ssh_stns_wrapper"); err != nil {
		fmt.Print(err)
		return
	}

	config, err := libnss_stns.LoadConfig(configFile)
	if err != nil {
		log.Print(err)
		return
	}

	if keys := FetchKey(flag.Arg(0), config); keys != "" {
		fmt.Println(keys)
	}
}

func FetchKey(name string, config *libnss_stns.Config) string {

	s := []string{config.ApiEndPoint, "user", "name", name}

	attr, err := libnss_stns.Request(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return ""
	}

	if attr != nil {
		return strings.Join(attr.Keys, "\n")
	}
	return ""
}

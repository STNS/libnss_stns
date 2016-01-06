package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/pyama86/libnss_stns/internal"
)

func main() {
	flag.Parse()
	if err := libnss_stns.InitLogger("ssh_stns_wrapper"); err != nil {
		fmt.Print(err)
		return
	}

	if keys := FetchKey(flag.Arg(0), "/etc/stns/libnss_stns.conf"); keys != "" {
		fmt.Println(keys)
	}
}

func FetchKey(name string, configFile string) string {
	config, err := libnss_stns.LoadConfig(configFile)
	if err != nil {
		log.Print(err)
		return ""
	}

	s := []string{config.Api_End_Point, "user", "name", name}

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

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

	config, err := libnss_stns.LoadConfig()
	if err != nil {
		log.Print(err)
		return
	}
	s := []string{config.Api_End_Point, "user", "name", flag.Arg(0)}
	attr, err := libnss_stns.Request(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return
	}

	if attr != nil {
		fmt.Println(strings.Join(attr.Keys, "\n"))
	}
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/STNS/lib-stns/config"
	"github.com/STNS/lib-stns/logger"
	"github.com/STNS/lib-stns/request"
)

func main() {
	logger.Setlog()
	config, err := config.Load("/etc/stns/lib_stns.conf")
	if err == nil {
		flag.Parse()
		if raw := Fetch(config, flag.Arg(0)); raw != "" {
			fmt.Println(raw)
		}
	}
}

func Fetch(config *config.Config, path string) string {
	r, err := request.NewRequest(config, path)

	if err != nil {
		log.Fatal(err)
	}

	result, err := r.GetRaw()
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

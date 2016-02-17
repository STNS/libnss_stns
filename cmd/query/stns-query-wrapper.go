package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/STNS/libnss_stns/cmd"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/request"
)

func main() {
	config, err := cmd.LoadConfig()
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

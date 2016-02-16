package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/STNS/libnss_stns/request"
)

func main() {
	flag.Parse()
	if raw := Fetch(flag.Arg(0)); raw != "" {
		fmt.Println(raw)
	}
}

func Fetch(path string) string {
	r, err := request.NewRequest(path)

	if err != nil {
		log.Print(err)
		return ""
	}

	result, err := r.GetRaw()
	if err != nil {
		return ""
	}

	if err != nil {
		return ""
	}
	return string(result)
}

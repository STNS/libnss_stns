package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/pyama86/libnss_stns/internal"
)

func main() {
	flag.Parse()
	config := libnss_stns.LoadConfig()
	s := []string{config.Api_End_Point, "user", "name", flag.Arg(0)}
	attr := libnss_stns.Request(strings.Join(s, "/"))
	if attr != nil {
		fmt.Println(strings.Join(attr.Keys, "\n"))
	}
}

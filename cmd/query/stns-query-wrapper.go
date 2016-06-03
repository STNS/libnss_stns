package main

import (
	"flag"
	"os"

	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/logger"
	"github.com/STNS/libnss_stns/request"
)

func main() {
	logger.Setlog()
	config, err := config.Load("/etc/stns/libnss_stns.conf")
	if err == nil {
		flag.Parse()
		raw, err := Fetch(config, flag.Arg(0))
		if err == nil {
			os.Stdout.Write([]byte(raw + "\n"))
		} else {
			os.Stderr.Write([]byte(err.Error() + "\n"))
		}
	}
}

func Fetch(config *config.Config, path string) (string, error) {
	r, err := request.NewRequest(config, path)

	if err != nil {
		return "", err
	}

	result, err := r.GetRawData()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

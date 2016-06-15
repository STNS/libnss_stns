package main

import (
	"flag"
	"os"

	"github.com/STNS/libnss_stns/libstns"
)

func main() {
	libstns.Setlog()
	config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
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

func Fetch(config *libstns.Config, path string) (string, error) {
	r, err := libstns.NewRequest(config, path)

	if err != nil {
		return "", err
	}

	result, err := r.GetRawData()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

package cmd

import (
	"log"

	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/logger"
)

func LoadConfig() (*config.Config, error) {
	logger.Setlog()
	config, err := config.Load("/etc/stns/libnss_stns.conf")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return config, nil
}

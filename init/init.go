package libnss_stns

import (
	"fmt"
	"log"

	"github.com/pyama86/libnss_stns/config"
	"github.com/pyama86/libnss_stns/logger"
)

const configFile = "/etc/stns/libnss_stns.conf"

func Init() (*config.Config, error) {
	if err := logger.Init("libnss_stns"); err != nil {
		fmt.Print(err)
		return nil, err
	}
	config, err := config.Load(configFile)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return config, nil
}

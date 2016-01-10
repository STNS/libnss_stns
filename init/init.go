package libnss_stns

import (
	"fmt"
	"log"
	"log/syslog"
	"reflect"

	"github.com/pyama86/libnss_stns/config"
)

const configFile = "/etc/stns/libnss_stns.conf"

var Loaded *config.Config

func Init(name string) (*config.Config, error) {

	if reflect.ValueOf(Loaded).IsNil() {
		logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, name)
		if err != nil {
			// syslog not found
			fmt.Print(err)
		} else {
			log.SetOutput(logger)
		}

		config, err := config.Load(configFile)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		Loaded = config
	}
	return Loaded, nil
}

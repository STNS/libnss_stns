package libnss_stns

import (
	"log"
	"log/syslog"
)

func InitLogger(name string) error {
	logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, name)
	if err != nil {
		return err
	}
	log.SetOutput(logger)
	return nil
}

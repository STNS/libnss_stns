package logger

import (
	"log"
	"log/syslog"
)

func Init(name string) error {
	logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, name)
	if err != nil {
		return err
	}
	log.SetOutput(logger)
	return nil
}

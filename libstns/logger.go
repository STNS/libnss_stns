package libstns

import (
	"log"
	"log/syslog"
)

func Setlog() error {
	logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "libstns")
	if err == nil {
		log.SetOutput(logger)
	}
	return nil
}

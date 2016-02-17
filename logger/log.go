package logger

import (
	"log"
	"log/syslog"
	"os"
)

func Setlog() {
	logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, os.Args[0])
	if err != nil {
		// syslog not found
		os.Stderr.Write([]byte(err))
	} else {
		log.SetOutput(logger)
	}
}

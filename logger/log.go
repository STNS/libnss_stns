package logger

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

func Setlog() {
	logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, os.Args[0])
	if err != nil {
		// syslog not found
		fmt.Print(err)
	} else {
		log.SetOutput(logger)
	}
}

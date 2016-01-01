package main

/*
#include <pwd.h>
#include <sys/types.h>
*/
import "C"
import (
	"log"
	"log/syslog"

	app "github.com/pyama86/libnss_etcd"
)

//export _nss_etcd_getpwnam_r
func _nss_etcd_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	config := app.LoadConfig()
	logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "my-daemon")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)

	log.Println("call name:", C.GoString(name))
	return 0
}

//export _nss_etcd_getspnam_r
func _nss_etcd_getspnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	config := app.LoadConfig()
	logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "my-daemon")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)

	log.Println("call name:", C.GoString(name))
	return 0
}

//export _nss_etcd_getpwuid_r
func _nss_etcd_getpwuid_r(uid C.uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	config := app.LoadConfig()
	logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "my-daemon")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logger)

	log.Println("call uid:", uid)
	return 0
}
func main() {
	config := app.LoadConfig()
}

package main

import (
	"log"
	"strconv"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/request"
)

/*
#include <pwd.h>
#include <sys/types.h>
*/
import "C"

var passwdList = map[string]*attribute.All{}
var passwdReadPos int

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	status := GetPasswd(pwd, result, "name", C.GoString(name))
	return status
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	status := GetPasswd(pwd, result, "id", strconv.Itoa(int(uid)))
	return status
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() {
	setList("user", passwdList, &passwdReadPos)
}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	resetList(passwdList, &passwdReadPos)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	passwds := getNextResource(passwdList, &passwdReadPos)
	if len(passwds) > 0 {
		setPasswd(pwd, passwds)
		result = &pwd
		return 1
	}
	return 0
}

func GetPasswd(pwd *C.struct_passwd, result **C.struct_passwd, column string, value string) int {
	passwds, err := request.Get("user", column, value)
	if err != nil {
		log.Print(err)
		return 0
	}
	if len(passwds) > 0 {
		setPasswd(pwd, passwds)
		result = &pwd
		return 1
	}
	return 0
}

func setPasswd(pwd *C.struct_passwd, passwds map[string]*attribute.All) {
	for n, p := range passwds {
		pwd.pw_name = C.CString(n)
		pwd.pw_passwd = C.CString("x")
		pwd.pw_uid = C.__uid_t(p.Id)
		pwd.pw_gid = C.__gid_t(p.GroupId)
		pwd.pw_gecos = C.CString(p.Gecos)
		pwd.pw_dir = C.CString(p.Directory)
		pwd.pw_shell = C.CString(p.Shell)
		return
	}
}

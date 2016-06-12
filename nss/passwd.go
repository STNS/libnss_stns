package main

import (
	"reflect"
	"strconv"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

/*
#include <pwd.h>
#include <sys/types.h>
*/
import "C"

var passwdList = stns.Attributes{}
var passwdReadPos int

func setPasswd(passwds stns.Attributes, p, r interface{}) int {
	pwd := p.(*C.struct_passwd)

	for n, p := range passwds {
		if p.User != nil && !reflect.ValueOf(p.User).IsNil() {
			pwd.pw_uid = C.__uid_t(p.Id)
			pwd.pw_name = C.CString(n)

			dir := "/home/" + n
			shell := "/bin/bash"

			if p.Directory != "" {
				dir = p.Directory
			}

			if p.Shell != "" {
				shell = p.Shell
			}
			pwd.pw_gid = C.__gid_t(p.GroupId)
			pwd.pw_passwd = C.CString("x")
			pwd.pw_dir = C.CString(dir)
			pwd.pw_shell = C.CString(shell)
			pwd.pw_gecos = C.CString(p.Gecos)
			r = &pwd
			return libstns.NSS_STATUS_SUCCESS
		}
	}
	return libstns.NSS_STATUS_NOTFOUND
}

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) C.int {
	return set(setPasswd, pwd, result, "user", "name", C.GoString(name))
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) C.int {
	return set(setPasswd, pwd, result, "user", "id", strconv.Itoa(int(uid)))
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() C.int {
	return initList(passwdList, &passwdReadPos, "user", libstns.NSS_LIST_PRESET)

}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	initList(passwdList, &passwdReadPos, "user", libstns.NSS_LIST_PURGE)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) C.int {
	return setByList(setPasswd, pwd, result, passwdList, &passwdReadPos)
}

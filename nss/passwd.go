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

type Passwd struct {
	pwd    *C.struct_passwd
	result **C.struct_passwd
}

func (s Passwd) Set(passwds stns.Attributes) int {
	for n, p := range passwds {
		if p.User != nil && !reflect.ValueOf(p.User).IsNil() {
			s.pwd.pw_uid = C.__uid_t(p.Id)
			s.pwd.pw_name = C.CString(n)

			dir := "/home/" + n
			shell := "/bin/bash"

			if p.Directory != "" {
				dir = p.Directory
			}

			if p.Shell != "" {
				shell = p.Shell
			}
			s.pwd.pw_gid = C.__gid_t(p.GroupId)
			s.pwd.pw_passwd = C.CString("x")
			s.pwd.pw_dir = C.CString(dir)
			s.pwd.pw_shell = C.CString(shell)
			s.pwd.pw_gecos = C.CString(p.Gecos)
			s.result = &s.pwd
			return libstns.NSS_STATUS_SUCCESS
		}
	}
	return libstns.NSS_STATUS_NOTFOUND
}

//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) C.int {
	p := Passwd{pwd, result}
	return set(pwdNss, p, "name", C.GoString(name))
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) C.int {
	p := Passwd{pwd, result}
	return set(pwdNss, p, "id", strconv.Itoa(int(uid)))
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() C.int {
	return initList(pwdNss, libstns.NSS_LIST_PRESET)

}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	initList(pwdNss, libstns.NSS_LIST_PURGE)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) C.int {
	p := Passwd{pwd, result}
	return setByList(pwdNss, p)
}

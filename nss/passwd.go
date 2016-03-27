package main

import (
	"reflect"
	"strconv"

	"github.com/STNS/STNS/attribute"
)

/*
#include <pwd.h>
#include <sys/types.h>
*/
import "C"

var passwdList = attribute.AllAttribute{}
var passwdReadPos int

type Passwd struct {
	pwd    *C.struct_passwd
	result **C.struct_passwd
}

func (self Passwd) setCStruct(passwds attribute.AllAttribute) int {

	for n, p := range passwds {
		if p.User != nil && !reflect.ValueOf(p.User).IsNil() {
			self.pwd.pw_uid = C.__uid_t(p.Id)
			self.pwd.pw_name = C.CString(n)

			dir := "/home/" + n
			shell := "/bin/bash"

			if p.Directory != "" {
				dir = p.Directory
			}

			if p.Shell != "" {
				shell = p.Shell
			}
			self.pwd.pw_gid = C.__gid_t(p.GroupId)
			self.pwd.pw_passwd = C.CString("x")
			self.pwd.pw_dir = C.CString(dir)
			self.pwd.pw_shell = C.CString(shell)
			self.pwd.pw_gecos = C.CString(p.Gecos)
			*self.result = self.pwd
			return NSS_STATUS_SUCCESS
		}
	}
	return NSS_STATUS_NOTFOUND
}

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setResource(&Passwd{pwd, result}, "user", "name", C.GoString(name))
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setResource(&Passwd{pwd, result}, "user", "id", strconv.Itoa(int(uid)))
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() int {
	return setList("user", passwdList, &passwdReadPos)

}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	resetList(passwdList, &passwdReadPos)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setNextResource(&Passwd{pwd, result}, passwdList, &passwdReadPos)
}

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

var passwdList = attribute.UserGroups{}
var passwdReadPos int

type Passwd struct {
	pwd    *C.struct_passwd
	result **C.struct_passwd
}

func (self Passwd) setCStruct(passwds attribute.UserGroups) {
	var dir, shell, gecos string
	var gid int

	for n, p := range passwds {
		self.pwd.pw_name = C.CString(n)

		dir = "/home/" + n
		shell = "/bin/bash"

		if p.User != nil && !reflect.ValueOf(p.User).IsNil() {
			self.pwd.pw_uid = C.__uid_t(p.Id)

			if p.Directory != "" {
				dir = p.Directory
			}

			if p.Shell != "" {
				shell = p.Shell
			}
			gid = p.GroupId
			gecos = p.GeCos
		}
		self.pwd.pw_gid = C.__gid_t(gid)
		self.pwd.pw_passwd = C.CString("x")
		self.pwd.pw_dir = C.CString(dir)
		self.pwd.pw_shell = C.CString(shell)
		self.pwd.pw_gecos = C.CString(gecos)

		*self.result = self.pwd
		return
	}
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
func _nss_stns_setpwent() {
	setList("user", passwdList, &passwdReadPos)

}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	resetList(passwdList, &passwdReadPos)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setNextResource(&Passwd{pwd, result}, passwdList, &passwdReadPos)
}

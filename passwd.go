package main

import (
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
	for n, p := range passwds {
		dir := "/home/" + n
		shell := "/bin/bash"

		if p.Directory != "" {
			dir = p.Directory
		}

		if p.Shell != "" {
			shell = p.Shell
		}
		self.pwd.pw_name = C.CString(n)
		self.pwd.pw_passwd = C.CString("x")
		self.pwd.pw_uid = C.__uid_t(p.Id)
		self.pwd.pw_gid = C.__gid_t(p.GroupId)
		self.pwd.pw_gecos = C.CString(p.Gecos)
		self.pwd.pw_dir = C.CString(dir)
		self.pwd.pw_shell = C.CString(shell)
		self.result = &self.pwd
		return
	}
}

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	r := Resource{"user"}
	return r.setResource(&Passwd{pwd, result}, "name", C.GoString(name))
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	r := Resource{"user"}
	return r.setResource(&Passwd{pwd, result}, "id", strconv.Itoa(int(uid)))
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() {
	entry := EntryResource{&Resource{"user"}, passwdList, &passwdReadPos}
	entry.setList()

}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	entry := EntryResource{&Resource{"user"}, passwdList, &passwdReadPos}
	entry.resetList()
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	entry := EntryResource{&Resource{"user"}, passwdList, &passwdReadPos}
	return entry.setNextResource(&Passwd{pwd, result})
}

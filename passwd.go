package main

import (
	"strconv"

	"github.com/pyama86/STNS/attribute"
)

/*
#include <pwd.h>
#include <sys/types.h>
*/
import "C"

var passwdList = attribute.UserGroups{}
var passwdReadPos int

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return getResource("user", "name", C.GoString(name), pwd, result)
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return getResource("user", "id", strconv.Itoa(int(uid)), pwd, result)
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() {
	setRecursiveList("user", passwdList, &passwdReadPos)
}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	resetRecursiveList(passwdList, &passwdReadPos)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return getRecursiveResource(pwd, result, passwdList, &passwdReadPos)
}

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

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setResource("user", "name", C.GoString(name), pwd, result)
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setResource("user", "id", strconv.Itoa(int(uid)), pwd, result)
}

//export _nss_stns_setpwent
func _nss_stns_setpwent() {
	setResourcePool("user", passwdList, &passwdReadPos)
}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	resetResourcePool(passwdList, &passwdReadPos)
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	return setResourceByPool(pwd, result, passwdList, &passwdReadPos)
}

package main

import "github.com/pyama86/STNS/attribute"

/*
#include <shadow.h>
*/
import "C"

var shadowList = attribute.UserGroups{}
var shadowReadPos int

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	return setResource("user", "name", C.GoString(name), spwd, result)
}

//export _nss_stns_setspent
func _nss_stns_setspent() {
	setRecursiveList("user", shadowList, &shadowReadPos)
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	resetRecursiveList(shadowList, &shadowReadPos)
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	return setRecursiveResource(spwd, result, shadowList, &shadowReadPos)
}

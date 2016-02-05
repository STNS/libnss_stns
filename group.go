package main

import (
	"strconv"

	"github.com/STNS/STNS/attribute"
)

/*
#include <grp.h>
#include <sys/types.h>
*/
import "C"

var groupList = attribute.UserGroups{}
var groupReadPos int

/*-------------------------------------------------------
group
-------------------------------------------------------*/
//export _nss_stns_getgrnam_r
func _nss_stns_getgrnam_r(name *C.char, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	return setResource("group", "name", C.GoString(name), grp, result)
}

//export _nss_stns_getgrgid_r
func _nss_stns_getgrgid_r(gid C.__gid_t, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	return setResource("group", "id", strconv.Itoa(int(gid)), grp, result)
}

//export _nss_stns_setgrent
func _nss_stns_setgrent() {
	setResourcePool("group", groupList, &groupReadPos)
}

//export _nss_stns_getgrent_r
func _nss_stns_getgrent_r(grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	return setResourceByPool(grp, result, groupList, &groupReadPos)
}

//export _nss_stns_endgrent
func _nss_stns_endgrent() {
	resetResourcePool(groupList, &groupReadPos)
}

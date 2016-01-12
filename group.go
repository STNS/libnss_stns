package main

import (
	"log"
	"strconv"
	"unsafe"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/request"
)

/*
#include <grp.h>
#include <sys/types.h>
*/
import "C"

var groupList = map[string]*attribute.All{}
var groupReadPos int

/*-------------------------------------------------------
group
-------------------------------------------------------*/
//export _nss_stns_getgrnam_r
func _nss_stns_getgrnam_r(name *C.char, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	status := GetGroup(grp, result, "name", C.GoString(name))
	return status
}

//export _nss_stns_getgrgid_r
func _nss_stns_getgrgid_r(gid C.__gid_t, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	status := GetGroup(grp, result, "id", strconv.Itoa(int(gid)))
	return status
}

//export _nss_stns_getgrent_r
func _nss_stns_getgrent_r(grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	groups := getNextResource(groupList, &groupReadPos)

	if len(groups) > 0 {
		setGroup(grp, groups)
		result = &grp
		return 1
	}
	return 0
}

//export _nss_stns_setgrent
func _nss_stns_setgrent() {
	setList("group", groupList, &groupReadPos)
}

//export _nss_stns_endgrent
func _nss_stns_endgrent() {
	resetList(groupList, &groupReadPos)
}

func GetGroup(grp *C.struct_group, result **C.struct_group, column string, value string) int {

	groups, err := request.Get("group", column, value)
	if err != nil {
		log.Print(err)
		return 0
	}

	if len(groups) > 0 {
		setGroup(grp, groups)
		result = &grp
		return 1
	}
	return 0
}

func setGroup(grp *C.struct_group, groups map[string]*attribute.All) {
	for n, g := range groups {
		grp.gr_name = C.CString(n)
		grp.gr_passwd = C.CString("x")
		grp.gr_gid = C.__gid_t(g.Id)
		work := make([]*C.char, len(g.Users)+1)
		if len(g.Users) > 0 {
			for i, u := range g.Users {
				work[i] = C.CString(u)
			}
		}
		grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
		return
	}
}

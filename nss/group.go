package main

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

/*
#include <grp.h>
#include <sys/types.h>
*/
import "C"

var groupList = stns.Attributes{}
var groupReadPos int

func setGroup(groups stns.Attributes, g, r interface{}) int {
	grp := g.(*C.struct_group)

	for n, g := range groups {
		if g.Id != 0 {
			grp.gr_gid = C.__gid_t(g.Id)
			grp.gr_name = C.CString(n)
			grp.gr_passwd = C.CString("x")

			if g.Group != nil && !reflect.ValueOf(g.Group).IsNil() {
				work := make([]*C.char, len(g.Users)+1)
				if len(g.Users) > 0 {
					for i, u := range g.Users {
						work[i] = C.CString(u)
					}
				}
				grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
			} else {
				work := make([]*C.char, 1)
				grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
			}

			r = &grp
			return libstns.NSS_STATUS_SUCCESS
		}
	}
	return libstns.NSS_STATUS_NOTFOUND
}

//export _nss_stns_getgrnam_r
func _nss_stns_getgrnam_r(name *C.char, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) C.int {
	return set(setGroup, grp, result, "group", "name", C.GoString(name))
}

//export _nss_stns_getgrgid_r
func _nss_stns_getgrgid_r(gid C.__gid_t, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) C.int {
	return set(setGroup, grp, result, "group", "id", strconv.Itoa(int(gid)))
}

//export _nss_stns_setgrent
func _nss_stns_setgrent() C.int {
	return initList(groupList, &groupReadPos, "group", libstns.NSS_LIST_PRESET)
}

//export _nss_stns_endgrent
func _nss_stns_endgrent() {
	initList(groupList, &groupReadPos, "group", libstns.NSS_LIST_PURGE)
}

//export _nss_stns_getgrent_r
func _nss_stns_getgrent_r(grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) C.int {
	return setByList(setGroup, grp, result, groupList, &groupReadPos)
}

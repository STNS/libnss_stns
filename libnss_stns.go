package main

/*
#include <pwd.h>
#include <grp.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
*/
import "C"
import (
	"strconv"
	"strings"
	"unsafe"

	"github.com/pyama86/libnss_stns/internal"
)

/*
 user
*/

//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	status := getPasswd(pwd, result, "user", "name", C.GoString(name))
	return status
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	status := getPasswd(pwd, result, "user", "id", strconv.Itoa(int(uid)))
	return status
}

func getPasswd(pwd *C.struct_passwd, result **C.struct_passwd, resource string, column string, value string) int {
	config := libnss_stns.LoadConfig()
	s := []string{config.Api_End_Point, resource, column, value}

	attr := libnss_stns.Request(strings.Join(s, "/"))

	if attr != nil {
		pwd.pw_name = C.CString(attr.Name)
		pwd.pw_passwd = C.CString("x")
		pwd.pw_uid = C.__uid_t(attr.Id)
		pwd.pw_gid = C.__gid_t(attr.Group_Id)
		pwd.pw_gecos = C.CString(attr.Gecos)
		pwd.pw_dir = C.CString(attr.Directory)
		pwd.pw_shell = C.CString(attr.Shell)
		result = &pwd
		return 1
	} else {
		return 0
	}
}

/*
 group
*/
//export _nss_stns_getgrnam_r
func _nss_stns_getgrnam_r(name *C.char, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	status := getGroup(grp, result, "group", "name", C.GoString(name))
	return status
}

//export _nss_stns_getgrgid_r
func _nss_stns_getgrgid_r(gid C.__gid_t, grp *C.struct_group, buffer *C.char, bufsize C.size_t, result **C.struct_group) int {
	status := getGroup(grp, result, "group", "id", strconv.Itoa(int(gid)))
	return status
}

func getGroup(grp *C.struct_group, result **C.struct_group, resource string, column string, value string) int {
	config := libnss_stns.LoadConfig()
	s := []string{config.Api_End_Point, resource, column, value}

	attr := libnss_stns.Request(strings.Join(s, "/"))
	if attr != nil {
		grp.gr_name = C.CString(attr.Name)
		grp.gr_passwd = C.CString("x")
		grp.gr_gid = C.__gid_t(attr.Id)
		work := make([]*C.char, len(attr.Users)+1)
		for i, u := range attr.Users {
			work[i] = C.CString(u)
		}

		grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
		result = &grp
		return 1
	} else {
		return 0
	}
}

func main() {
}

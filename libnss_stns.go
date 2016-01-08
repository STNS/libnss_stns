package main

/*
#include <pwd.h>
#include <shadow.h>
#include <grp.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
*/
import "C"
import (
	"log"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/init"
	"github.com/pyama86/libnss_stns/request"
)

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
//export _nss_stns_getpwnam_r
func _nss_stns_getpwnam_r(name *C.char, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	status := GetPasswd(pwd, result, "name", C.GoString(name))
	return status
}

//export _nss_stns_getpwuid_r
func _nss_stns_getpwuid_r(uid C.__uid_t, pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	status := GetPasswd(pwd, result, "id", strconv.Itoa(int(uid)))
	return status
}
func GetPasswd(pwd *C.struct_passwd, result **C.struct_passwd, column string, value string) int {
	config, err := libnss_stns.Init()
	if err != nil {
		return 0
	}

	s := []string{config.ApiEndPoint, "user", column, value}

	passwds, err := request.Send(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return 0
	}
	if len(passwds) > 0 {
		for n, p := range passwds {
			pwd.pw_name = C.CString(n)
			pwd.pw_passwd = C.CString("x")
			pwd.pw_uid = C.__uid_t(p.Id)
			pwd.pw_gid = C.__gid_t(p.GroupId)
			pwd.pw_gecos = C.CString(p.Gecos)
			pwd.pw_dir = C.CString(p.Directory)
			pwd.pw_shell = C.CString(p.Shell)
			result = &pwd
			return 1
		}
	}
	return 0
}

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	status := GetShadow(spwd, result, "name", C.GoString(name))
	return status
}
func GetShadow(spwd *C.struct_spwd, result **C.struct_spwd, column string, value string) int {
	config, err := libnss_stns.Init()
	if err != nil {
		return 0
	}

	s := []string{config.ApiEndPoint, "user", column, value}

	shadows, err := request.Send(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return 0
	}
	if len(shadows) > 0 {
		for n, _ := range shadows {
			spwd.sp_namp = C.CString(n)
			spwd.sp_pwdp = C.CString("!!")
			spwd.sp_lstchg = -1
			spwd.sp_min = -1
			spwd.sp_max = -1
			spwd.sp_warn = -1
			spwd.sp_inact = -1
			spwd.sp_expire = -1

			result = &spwd
			return 1
		}
	}
	return 0
}

/*-------------------------------------------------------
group
-------------------------------------------------------*/
var groupList map[string]*attribute.All
var groupReadPos int

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
	keys := groupKeys(groupList)
	if len(keys) > groupReadPos && keys[groupReadPos] != "" {
		name := keys[groupReadPos]

		grp.gr_name = C.CString(name)
		grp.gr_passwd = C.CString("!!")
		grp.gr_gid = C.__gid_t(groupList[name].Id)
		work := make([]*C.char, len(groupList[name].Users)+1)
		for i, u := range groupList[name].Users {
			work[i] = C.CString(u)
		}

		grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
		result = &grp
		groupReadPos++
		return 1

	}
	return 0
}

//export _nss_stns_setgrent
func _nss_stns_setgrent() {
	groupReadPos = 0
	config, err := libnss_stns.Init()
	if err != nil {
		return
	}

	s := []string{config.ApiEndPoint, "group", "list"}

	groupList, err = request.Send(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return
	}
	return
}

//export _nss_stns_endgrent
func _nss_stns_endgrent() {
	groupReadPos = 0
	groupList = nil
	return
}

func groupKeys(m map[string]*attribute.All) []string {
	ks := []string{}
	for k, _ := range m {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

func GetGroup(grp *C.struct_group, result **C.struct_group, column string, value string) int {
	config, err := libnss_stns.Init()
	if err != nil {
		return 0
	}

	s := []string{config.ApiEndPoint, "group", column, value}

	groups, err := request.Send(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return 0
	}

	if len(groups) > 0 {
		for n, g := range groups {
			grp.gr_name = C.CString(n)
			grp.gr_passwd = C.CString("!!")
			grp.gr_gid = C.__gid_t(g.Id)
			work := make([]*C.char, len(g.Users)+1)
			for i, u := range g.Users {
				work[i] = C.CString(u)
			}

			grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
			result = &grp
			return 1
		}
	}
	return 0
}

func ListGroup(grp *C.struct_group, result **C.struct_group) int {
	config, err := libnss_stns.Init()
	if err != nil {
		return 0
	}

	s := []string{config.ApiEndPoint, "group", "list"}

	groups, err := request.Send(strings.Join(s, "/"))
	if err != nil {
		log.Print(err)
		return 0
	}

	if len(groups) > 0 {
		for n, g := range groups {
			grp.gr_name = C.CString(n)
			grp.gr_passwd = C.CString("!!")
			grp.gr_gid = C.__gid_t(g.Id)
			work := make([]*C.char, len(g.Users)+1)
			for i, u := range g.Users {
				work[i] = C.CString(u)
			}
			grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
		}
		result = &grp
		return 1
	}
	return 0
}

func main() {
}

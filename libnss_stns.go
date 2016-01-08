package main

/*
#include <pwd.h>
#include <shadow.h>
#include <grp.h>
#include <sys/types.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
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

func getList(resource string) (map[string]*attribute.All, error) {
	config, err := libnss_stns.Init()
	if err != nil {
		return nil, err
	}
	s := []string{config.ApiEndPoint, resource, "list"}

	list, err := request.Send(strings.Join(s, "/"))

	if err != nil {
		return nil, err
	}
	return list, err
}

func getKeys(m map[string]*attribute.All) []string {
	ks := []string{}
	for k, _ := range m {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

func getNextResource(list map[string]*attribute.All, pos *int) (string, *attribute.All) {
	keys := getKeys(list)
	if len(keys) > *pos && keys[*pos] != "" {
		name := keys[*pos]
		resource := list[name]
		*pos++
		return name, resource
	}
	return "", nil
}

/*-------------------------------------------------------
passwd
-------------------------------------------------------*/
var passwdList map[string]*attribute.All
var passwdReadPos int

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

//export _nss_stns_setpwent
func _nss_stns_setpwent() {
	var err error
	passwdReadPos = 0
	passwdList, err = getList("user")
	if err != nil {
		log.Print(err)
	}
}

//export _nss_stns_endpwent
func _nss_stns_endpwent() {
	passwdList = nil
	passwdReadPos = 0
}

//export _nss_stns_getpwent_r
func _nss_stns_getpwent_r(pwd *C.struct_passwd, buffer *C.char, bufsize C.size_t, result **C.struct_passwd) int {
	name, passwd := getNextResource(passwdList, &passwdReadPos)
	if name != "" {
		pwd.pw_name = C.CString(name)
		pwd.pw_passwd = C.CString("x")
		pwd.pw_uid = C.__uid_t(passwd.Id)
		pwd.pw_gid = C.__gid_t(passwd.GroupId)
		pwd.pw_gecos = C.CString(passwd.Gecos)
		pwd.pw_dir = C.CString(passwd.Directory)
		pwd.pw_shell = C.CString(passwd.Shell)
		result = &pwd
		return 1
	}
	return 0
}

func GetPasswd(pwd *C.struct_passwd, result **C.struct_passwd, column string, value string) int {
	passwds, err := getList("user")
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
			*result = pwd
			return 1
		}
	}
	return 0
}

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/
var shadowList map[string]*attribute.All
var shadowReadPos int

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	status := GetShadow(spwd, result, "name", C.GoString(name))
	return status
}

//export _nss_stns_setspent
func _nss_stns_setspent() {
	var err error
	shadowReadPos = 0
	shadowList, err = getList("user")
	if err != nil {
		log.Print(err)
	}
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	shadowList = nil
	shadowReadPos = 0
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	name, _ := getNextResource(shadowList, &shadowReadPos)
	if name != "" {
		spwd.sp_namp = C.CString(name)
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
	return 0
}

func GetShadow(spwd *C.struct_spwd, result **C.struct_spwd, column string, value string) int {
	shadows, err := getList("user")
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
	name, group := getNextResource(groupList, &groupReadPos)
	if name != "" {
		grp.gr_name = C.CString(name)
		grp.gr_passwd = C.CString("!!")
		grp.gr_gid = C.__gid_t(group.Id)
		work := make([]*C.char, len(group.Users)+1)
		for i, u := range group.Users {
			work[i] = C.CString(u)
		}

		grp.gr_mem = (**C.char)(unsafe.Pointer(&work[0]))
		result = &grp
		return 1
	}
	return 0
}

//export _nss_stns_setgrent
func _nss_stns_setgrent() {
	var err error
	groupReadPos = 0
	groupList, err = getList("group")
	if err != nil {
		log.Print(err)
	}
}

//export _nss_stns_endgrent
func _nss_stns_endgrent() {
	groupList = nil
	groupReadPos = 0
}

func GetGroup(grp *C.struct_group, result **C.struct_group, column string, value string) int {

	groups, err := getList("group")
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

func main() {
}

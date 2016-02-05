package main

import (
	"log"
	"sort"
	"unsafe"

	"github.com/STNS/STNS/attribute"
	"github.com/STNS/libnss_stns/request"
)

/*
#include <grp.h>
#include <pwd.h>
#include <shadow.h>
#include <sys/types.h>
*/
import "C"

func setResource(resource_name string, column string, value string, obj interface{}, result interface{}) int {
	r, err := request.NewRequest(resource_name, column, value)
	if err != nil {
		log.Print(err)
		return -2
	}

	resource, err := r.Get()
	if err != nil {
		log.Print(err)
		return -2
	}
	if len(resource) > 0 {
		setLinuxResource(obj, result, resource)
		return 1
	}
	return 0
}

func setRecursiveResource(obj interface{}, result interface{}, list attribute.UserGroups, pos *int) int {
	keys := getKeys(list)
	if len(keys) > *pos && keys[*pos] != "" {
		name := keys[*pos]
		resource := attribute.UserGroups{
			name: list[name],
		}

		setLinuxResource(obj, result, resource)
		*pos++
		return 1
	}
	return 0
}

func setLinuxResource(obj interface{}, result interface{}, resource attribute.UserGroups) {
	switch obj.(type) {
	case *C.struct_passwd:
		setPasswd(obj.(*C.struct_passwd), resource)
		result = &obj
	case *C.struct_group:
		setGroup(obj.(*C.struct_group), resource)
		result = &obj
	case *C.struct_spwd:
		setShadow(obj.(*C.struct_spwd), resource)
		result = &obj
	case attribute.UserGroups:
		for k, v := range resource {
			obj.(attribute.UserGroups)[k] = v
		}
	}
}

func setPasswd(pwd *C.struct_passwd, passwds attribute.UserGroups) {
	for n, p := range passwds {
		dir := "/home/" + n
		shell := "/bin/bash"

		if p.Directory != "" {
			dir = p.Directory
		}

		if p.Shell != "" {
			shell = p.Shell
		}
		pwd.pw_name = C.CString(n)
		pwd.pw_passwd = C.CString("x")
		pwd.pw_uid = C.__uid_t(p.Id)
		pwd.pw_gid = C.__gid_t(p.GroupId)
		pwd.pw_gecos = C.CString(p.Gecos)
		pwd.pw_dir = C.CString(dir)
		pwd.pw_shell = C.CString(shell)
		return
	}
}

func setShadow(spwd *C.struct_spwd, shadows attribute.UserGroups) {
	for n, _ := range shadows {
		spwd.sp_namp = C.CString(n)
		spwd.sp_pwdp = C.CString("!!")
		spwd.sp_lstchg = -1
		spwd.sp_min = -1
		spwd.sp_max = -1
		spwd.sp_warn = -1
		spwd.sp_inact = -1
		spwd.sp_expire = -1
		return
	}
}
func setGroup(grp *C.struct_group, groups attribute.UserGroups) {
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

func setRecursiveList(resource_name string, list attribute.UserGroups, pos *int) {
	// reset value
	resetRecursiveList(list, pos)
	err := setResource(resource_name, "list", "", list, nil)
	if err != 1 {
		log.Print(err)
		return
	}
}

func resetRecursiveList(list attribute.UserGroups, pos *int) {
	// reset value
	*pos = 0
	for k, _ := range list {
		delete(list, k)
	}
}

func getKeys(m attribute.UserGroups) []string {
	ks := []string{}
	for k, _ := range m {
		ks = append(ks, k)

	}
	sort.Strings(ks)
	return ks
}

package main

import (
	"log"

	"github.com/pyama86/STNS/attribute"
	"github.com/pyama86/libnss_stns/request"
)

/*
#include <shadow.h>
*/
import "C"

var shadowList = map[string]*attribute.All{}
var shadowReadPos int

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	status := GetShadow(spwd, result, "name", C.GoString(name))
	return status
}

//export _nss_stns_setspent
func _nss_stns_setspent() {
	setList("user", shadowList, &shadowReadPos)
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	shadowList = nil
	shadowReadPos = 0
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	shadows := getNextResource(shadowList, &shadowReadPos)
	if len(shadows) > 0 {
		setShadow(spwd, shadows)
		result = &spwd
		return 1
	}
	return 0
}
func GetShadow(spwd *C.struct_spwd, result **C.struct_spwd, column string, value string) int {
	shadows, err := request.Get("user", column, value)
	if err != nil {
		log.Print(err)
		return 0
	}
	if len(shadows) > 0 {
		setShadow(spwd, shadows)
		result = &spwd
		return 1
	}
	return 0
}

func setShadow(spwd *C.struct_spwd, shadows map[string]*attribute.All) {
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

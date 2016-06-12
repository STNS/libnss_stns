package main

/*
#include <shadow.h>
*/
import "C"
import (
	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

var shadowList = stns.Attributes{}
var shadowReadPos int

func setShadow(shadows stns.Attributes, s, r interface{}) int {
	spwd := s.(*C.struct_spwd)

	for n, _ := range shadows {
		spwd.sp_namp = C.CString(n)
		spwd.sp_pwdp = C.CString("!!")
		spwd.sp_lstchg = -1
		spwd.sp_min = -1
		spwd.sp_max = -1
		spwd.sp_warn = -1
		spwd.sp_inact = -1
		spwd.sp_expire = -1
		r = &spwd
		return libstns.NSS_STATUS_SUCCESS
	}
	return libstns.NSS_STATUS_NOTFOUND
}

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) C.int {
	return set(setShadow, spwd, result, "user", "name", C.GoString(name))
}

//export _nss_stns_setspent
func _nss_stns_setspent() C.int {
	return initList(shadowList, &shadowReadPos, "user", libstns.NSS_LIST_PRESET)
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	initList(shadowList, &shadowReadPos, "user", libstns.NSS_LIST_PURGE)
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) C.int {
	return setByList(setShadow, spwd, result, shadowList, &shadowReadPos)
}

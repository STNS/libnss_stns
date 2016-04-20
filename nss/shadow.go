package main

/*
#include <shadow.h>
*/
import "C"
import "github.com/STNS/STNS/stns"

var shadowList = stns.Attributes{}
var shadowReadPos int

type Shadow struct {
	spwd   *C.struct_spwd
	result **C.struct_spwd
}

func (self Shadow) setCStruct(shadows stns.Attributes) int {
	for n, _ := range shadows {
		self.spwd.sp_namp = C.CString(n)
		self.spwd.sp_pwdp = C.CString("!!")
		self.spwd.sp_lstchg = -1
		self.spwd.sp_min = -1
		self.spwd.sp_max = -1
		self.spwd.sp_warn = -1
		self.spwd.sp_inact = -1
		self.spwd.sp_expire = -1
		self.result = &self.spwd
		return NSS_STATUS_SUCCESS
	}
	return NSS_STATUS_NOTFOUND
}

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	return set(&Shadow{spwd, result}, "user", "name", C.GoString(name))
}

//export _nss_stns_setspent
func _nss_stns_setspent() int {
	return initList("user", shadowList, &shadowReadPos)
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	purgeList(shadowList, &shadowReadPos)
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	return setByList(&Shadow{spwd, result}, shadowList, &shadowReadPos)
}

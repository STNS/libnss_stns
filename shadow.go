package main

import "github.com/STNS/STNS/attribute"

/*
#include <shadow.h>
*/
import "C"

var shadowList = attribute.UserGroups{}
var shadowReadPos int

type Shadow struct {
	spwd   *C.struct_spwd
	result **C.struct_spwd
}

func (self Shadow) setCStruct(shadows attribute.UserGroups) {
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
		return
	}
}

/*-------------------------------------------------------
shadow
-------------------------------------------------------*/

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	r := Resource{"user"}
	return r.setResource(&Shadow{spwd, result}, "name", C.GoString(name))
}

//export _nss_stns_setspent
func _nss_stns_setspent() {
	entry := EntryResource{&Resource{"user"}, shadowList, &shadowReadPos}
	entry.setList()
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	entry := EntryResource{&Resource{"user"}, shadowList, &shadowReadPos}
	entry.resetList()
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) int {
	entry := EntryResource{&Resource{"user"}, shadowList, &shadowReadPos}
	return entry.setNextResource(&Shadow{spwd, result})
}

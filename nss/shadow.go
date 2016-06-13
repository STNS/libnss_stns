package main

/*
#include <shadow.h>
*/
import "C"
import (
	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

type Shadow struct {
	spwd   *C.struct_spwd
	result **C.struct_spwd
}

func (s Shadow) Set(shadows stns.Attributes) int {
	for n, _ := range shadows {
		s.spwd.sp_namp = C.CString(n)
		s.spwd.sp_pwdp = C.CString("!!")
		s.spwd.sp_lstchg = -1
		s.spwd.sp_min = -1
		s.spwd.sp_max = -1
		s.spwd.sp_warn = -1
		s.spwd.sp_inact = -1
		s.spwd.sp_expire = -1
		s.result = &s.spwd
		return libstns.NSS_STATUS_SUCCESS
	}
	return libstns.NSS_STATUS_NOTFOUND
}

//export _nss_stns_getspnam_r
func _nss_stns_getspnam_r(name *C.char, spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) C.int {
	s := Shadow{spwd, result}
	return set(spwdNss, s, "name", C.GoString(name))
}

//export _nss_stns_setspent
func _nss_stns_setspent() C.int {
	return initList(spwdNss, libstns.NSS_LIST_PRESET)
}

//export _nss_stns_endspent
func _nss_stns_endspent() {
	initList(spwdNss, libstns.NSS_LIST_PURGE)
}

//export _nss_stns_getspent_r
func _nss_stns_getspent_r(spwd *C.struct_spwd, buffer *C.char, bufsize C.size_t, result **C.struct_spwd) C.int {
	s := Shadow{spwd, result}
	return setByList(spwdNss, s)
}

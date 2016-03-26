package main

/*
#include <stdlib.h>
#include <security/pam_appl.h>
#include <security/pam_ext.h>
#cgo LDFLAGS: -lpam
*/
import "C"

// The reason that separates this method, but in order to avoid a compile error
func getPassword(pamh *C.pam_handle_t, item C.int, password **C.char) bool {
	if C.pam_get_authtok(pamh, C.PAM_AUTHTOK, password, nil) != C.PAM_SUCCESS {
		return false
	}
	return true
}

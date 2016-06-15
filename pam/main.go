package main

/*
#include <stdlib.h>
#include <security/pam_appl.h>

*/
import "C"
import (
	"log"
	"unsafe"

	"github.com/STNS/libnss_stns/libstns"
)

func main() {
}

func init() {
	libstns.Setlog()
}

//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) C.int {
	config, err := libstns.LoadConfig("/etc/stns/libnss_stns.conf")
	if err != nil {
		log.Println(err)
		return C.PAM_AUTHINFO_UNAVAIL
	}

	pam := libstns.NewPam(config, int(argc), GoStrings(int(argc), argv))

	var user string

	if pam.AuthType == "user" {
		var cuser *C.char
		defer C.free(unsafe.Pointer(cuser))
		if getUserName(pamh, &cuser) {
			user = C.GoString(cuser)
		} else {
			return C.PAM_USER_UNKNOWN
		}
	} else if pam.AuthType == "sudo" {
		user = pam.SudoUser()
	} else {
		return C.PAM_USER_UNKNOWN
	}

	var cPassword *C.char
	defer C.free(unsafe.Pointer(cPassword))
	if !getPassword(pamh, &cPassword) {
		return C.PAM_AUTH_ERR
	}
	return C.int(pam.PasswordAuth(user, C.GoString(cPassword)))
}

//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) C.int {
	return C.PAM_SUCCESS
}

func GoStrings(length int, argv **C.char) []string {
	if length > 0 {
		tmpslice := (*[1 << 27]*C.char)(unsafe.Pointer(argv))[:length:length]
		gostrings := make([]string, length)
		for i, s := range tmpslice {
			gostrings[i] = C.GoString(s)
		}
		return gostrings
	}
	return nil
}

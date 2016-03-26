package main

/*
#include <security/pam_appl.h>
*/
import "C"
import "unsafe"

// export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) int {
	return int(C.PAM_AUTH_ERR)
}

//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) int {
	return int(C.PAM_SUCCESS)
}
func GoStrings(argc C.int, argv **C.char) []string {

	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}
func main() {
}

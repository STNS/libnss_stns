package main

/*
#include <security/pam_appl.h>
*/
import "C"
import (
	"crypto/sha256"
	"log"
	"unsafe"

	"github.com/STNS/libnss_stns/logger"
	"github.com/STNS/libnss_stns/request"
)

// export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char) int {
	logger.Setlog()
	c := int(argc)

	if c > 0 {
		a := GoStrings(argc, argv)
		switch a[0] {
		case "sudo":
			var passwd *C.String
			defer C.free(unsafe.Pointer(passwd))
			if len(a) > 1 {
				if C.pam_get_authtok(pamh, C.PAM_AUTHTOK, password, nil) != C.PAM_SUCCESS {
					return int(C.PAM_AUTH_ERR)
				}

				r, err := request.NewRequest(config, "sudo", "name", a[1])
				if err != nil {
					log.Println(err)
					return int(C.PAM_AUTH_ERR)
				}

				sudo, err := r.Get()
				if err != nil {
					log.Println(err)
					return int(C.PAM_AUTH_ERR)
				}

				if sha256.Sum256(C.GoString(password)) == sudo.Password {
					return int(C.PAM_AUTH_SUCCESS)
				}
			}
			return int(C.PAM_AUTH_ERR)
		default:
			return int(C.PAM_AUTH_ERR)
		}
	}
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

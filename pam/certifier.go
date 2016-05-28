package main

/*
#include <stdlib.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>
*/
import "C"
import (
	"unsafe"

	"github.com/STNS/libnss_stns/config"
)

type Certifier interface {
	Auth(Certifier) C.int
	userName() string
}

type Supplicant struct {
	authType string
	pamh     *C.pam_handle_t
	argc     int
	argv     []string
	config   *config.Config
}

type Sudo struct {
	*Supplicant
}

func NewCertifier(pamh *C.pam_handle_t, argc C.int, argv **C.char, config *config.Config) Certifier {
	gargc := int(argc)
	gargv := GoStrings(gargc, argv)
	if gargc > 0 {
		switch gargv[0] {
		case "sudo":
			return Sudo{&Supplicant{"sudo", pamh, gargc, gargv, config}}
		}
	}
	return Supplicant{"user", pamh, gargc, gargv, config}
}

func (s Supplicant) userName() string {
	var user *C.char
	defer C.free(unsafe.Pointer(user))
	if authUser(s.pamh, &user) {
		return C.GoString(user)
	}
	return ""
}

func (s Sudo) userName() string {
	if s.argc > 1 {
		return s.argv[1]
	}
	return ""
}

func (s Supplicant) Auth(certifier Certifier) C.int {
	var cPassword *C.char
	defer C.free(unsafe.Pointer(cPassword))

	user := certifier.userName()
	if user == "" {
		return C.PAM_USER_UNKNOWN
	}

	if !authPassword(s.pamh, &cPassword) {
		return C.PAM_AUTH_ERR
	}

	return C.int(checkPassword(s.config, s.authType, user, C.GoString(cPassword)))
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

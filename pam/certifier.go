package main

/*
#include <stdlib.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>
*/
import "C"
import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"strings"
	"unsafe"

	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/request"
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
	var password *C.char
	defer C.free(unsafe.Pointer(password))

	user := certifier.userName()
	if user == "" {
		return C.PAM_USER_UNKNOWN
	}

	if !authPassword(s.pamh, &password) {
		return C.PAM_AUTH_ERR
	}

	r, err := request.NewRequest(s.config, s.authType, "name", user)
	if err != nil {
		log.Println(err)
		return C.PAM_AUTH_ERR
	}

	attr, err := r.Get()
	if err != nil {
		log.Println(err)
		return C.PAM_AUTHINFO_UNAVAIL
	}

	if attr != nil {
		for _, s := range attr {
			var hash string
			switch s.HashType {
			case "sha512":
				hash = strings.ToLower(sha512Sum([]byte(C.GoString(password))))
			default:
				hash = strings.ToLower(sha256Sum([]byte(C.GoString(password))))
			}
			if hash == strings.ToLower(s.Password) {
				return C.PAM_SUCCESS
			}

		}
	}
	return C.PAM_AUTH_ERR
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

func sha256Sum(data []byte) string {
	bytes := sha256.Sum256(data)
	return hex.EncodeToString(bytes[:])
}

func sha512Sum(data []byte) string {
	bytes := sha512.Sum512(data)
	return hex.EncodeToString(bytes[:])
}

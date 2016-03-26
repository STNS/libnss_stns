package main

/*
#include <stdlib.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>
*/
import "C"
import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
	"unsafe"

	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/request"
)

type Certifier interface {
	Auth() C.int
}
type Sudo struct {
	pamh   *C.pam_handle_t
	argc   int
	argv   []string
	config *config.Config
}

func getCertifier(pamh *C.pam_handle_t, argc C.int, argv **C.char, config *config.Config) Certifier {
	gargc := int(argc)
	if gargc > 0 {
		gargv := GoStrings(gargc, argv)
		switch gargv[0] {
		case "sudo":
			return Sudo{pamh, gargc, gargv, config}
		}
	}
	return nil
}

func (s Sudo) Auth() C.int {
	if s.argc > 1 {
		var password *C.char
		defer C.free(unsafe.Pointer(password))
		if !getPassword(s.pamh, C.PAM_AUTHTOK, &password) {
			return C.PAM_AUTH_ERR
		}

		r, err := request.NewRequest(s.config, "sudo", "name", s.argv[1])
		if err != nil {
			log.Println(err)
			return C.PAM_AUTH_ERR
		}

		sudoers, err := r.Get()
		if err != nil {
			log.Println(err)
			return C.PAM_AUTH_ERR
		}

		if sudoers != nil {
			for _, s := range sudoers {
				if strings.ToLower(sha256Sum([]byte(C.GoString(password)))) == strings.ToLower(s.Password) {
					return C.PAM_SUCCESS
				}
			}
		}
	}
	return C.PAM_AUTH_ERR
}

func GoStrings(length int, argv **C.char) []string {
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

func sha256Sum(data []byte) string {
	bytes := sha256.Sum256(data)
	return hex.EncodeToString(bytes[:])
}

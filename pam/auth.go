package main

import (
	"log"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/request"
)

const (
	PAM_AUTH_ERR         = 7
	PAM_AUTHINFO_UNAVAIL = 9
	PAM_SUCCESS          = 0
)

func checkPassword(config *config.Config, authType string, user string, password string) int {
	var attr stns.Attribute
	var salt string

	r, err := request.NewRequest(config, authType, "name", user)
	if err != nil {
		log.Println(err)
		return PAM_AUTHINFO_UNAVAIL
	}

	res, err := r.GetByWrapperCmd()
	if err != nil {
		log.Println(err)
		return PAM_AUTHINFO_UNAVAIL
	}

	for _, a := range *res.Items {
		attr = *a
		break
	}

	var hashMethod HashMethod
	switch attr.HashType {
	case "sha512":
		hashMethod = sha512Sum
	default:
		hashMethod = sha256Sum
	}

	if res.MetaData.Salt {
		salt += hashMethod([]byte(user))
	}

	hash := strings.ToLower(hashMethod([]byte(salt + password)))

	for i := 0; i < res.MetaData.Stretching-1; i++ {
		hash = strings.ToLower(hashMethod([]byte(hash)))
	}

	r.SetPath("auth", authType, "name", user, hash)
	ar, err := r.GetByWrapperCmd()
	if err != nil {
		log.Println(err)
		return PAM_AUTHINFO_UNAVAIL
	}

	if ar.MetaData.Result == "success" {
		return PAM_SUCCESS
	}
	return PAM_AUTH_ERR

}

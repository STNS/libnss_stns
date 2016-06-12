package main

import (
	"log"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/libstns"
)

const (
	PAM_AUTH_ERR         = 7
	PAM_AUTHINFO_UNAVAIL = 9
	PAM_SUCCESS          = 0
)

func checkPassword(config *libstns.Config, authType string, user string, password string) int {
	var attr stns.Attribute
	var hashType string
	r, err := libstns.NewRequest(config, authType, "name", user)
	if err != nil {
		log.Println(err)
		return PAM_AUTHINFO_UNAVAIL
	}
	res, err := r.GetByWrapperCmd()
	if err != nil {
		log.Println(err)
		return PAM_AUTHINFO_UNAVAIL
	}

	if res.Items == nil {
		log.Printf("resource notfound %s/%s", authType, user)
		return PAM_AUTHINFO_UNAVAIL
	}

	for _, a := range *res.Items {
		attr = *a
		break
	}

	if attr.HashType == "" {
		hashType = res.MetaData.HashType
	} else {
		hashType = attr.HashType
	}

	if strings.ToLower(attr.Password) == libstns.Calculate(hashType, res.MetaData.Salt, user, password, res.MetaData.Stretching) {
		return PAM_SUCCESS
	}

	return PAM_AUTH_ERR
}

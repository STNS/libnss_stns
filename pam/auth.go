package main

import (
	"log"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/STNS/libnss_stns/config"
	"github.com/STNS/libnss_stns/hash"
	"github.com/STNS/libnss_stns/request"
)

const (
	PAM_AUTH_ERR         = 7
	PAM_AUTHINFO_UNAVAIL = 9
	PAM_SUCCESS          = 0
)

func checkPassword(config *config.Config, authType string, user string, password string) int {
	var attr stns.Attribute

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

	if strings.ToLower(attr.Password) == hash.Calculate(attr.HashType, res.MetaData.Salt, user, password, res.MetaData.Stretching) {
		return PAM_SUCCESS
	}

	return PAM_AUTH_ERR
}

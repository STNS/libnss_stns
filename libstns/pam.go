package libstns

import (
	"log"
	"strings"

	"github.com/STNS/STNS/stns"
	"github.com/tredoe/osutil/user/crypt"
	"github.com/tredoe/osutil/user/crypt/apr1_crypt"
	"github.com/tredoe/osutil/user/crypt/md5_crypt"
	"github.com/tredoe/osutil/user/crypt/sha256_crypt"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

const (
	PAM_AUTH_ERR         = 7
	PAM_AUTHINFO_UNAVAIL = 9
	PAM_SUCCESS          = 0
)

type Pam struct {
	config   *Config
	AuthType string
	argc     int
	argv     []string
}

func NewPam(config *Config, argc int, argv []string) *Pam {
	var u string
	u = "user"
	if argc > 0 {
		u = argv[0]
	}

	return &Pam{
		config:   config,
		AuthType: u,
		argc:     argc,
		argv:     argv,
	}
}

func (p *Pam) SudoUser() string {
	if p.argc > 1 {
		return p.argv[1]
	}
	return ""
}

func (p *Pam) PasswordAuth(user string, password string) int {
	r, err := NewRequest(p.config, p.AuthType, "name", user)
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
		log.Printf("resource notfound %s/%s", p.AuthType, user)
		return PAM_AUTHINFO_UNAVAIL
	}

	var attr stns.Attribute
	for _, a := range res.Items {
		attr = *a
		break
	}

	if strings.Count(attr.Password, "$") != 3 {
		return PAM_AUTHINFO_UNAVAIL
	}

	var c crypt.Crypter
	switch {
	case strings.HasPrefix(attr.Password, sha512_crypt.MagicPrefix):
		c = sha512_crypt.New()
	case strings.HasPrefix(attr.Password, sha256_crypt.MagicPrefix):
		c = sha256_crypt.New()
	case strings.HasPrefix(attr.Password, md5_crypt.MagicPrefix):
		c = md5_crypt.New()
	case strings.HasPrefix(attr.Password, apr1_crypt.MagicPrefix):
		c = apr1_crypt.New()
	}

	err = c.Verify(attr.Password, []byte(password))
	if err == nil {
		return PAM_SUCCESS
	}

	return PAM_AUTH_ERR
}

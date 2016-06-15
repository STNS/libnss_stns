package libstns

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"strings"

	"github.com/STNS/STNS/stns"
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
	for _, a := range *res.Items {
		attr = *a
		break
	}

	var hashType string
	hashType = attr.HashType

	if hashType == "" {
		hashType = res.MetaData.HashType
	}

	if strings.ToLower(attr.Password) == p.GenerateHash(hashType, user, password, res.MetaData.Salt, res.MetaData.Stretching) {
		return PAM_SUCCESS
	}

	return PAM_AUTH_ERR
}

type HashMethod func([]byte) string

func (p *Pam) sha256Sum(data []byte) string {
	bytes := sha256.Sum256(data)
	return hex.EncodeToString(bytes[:])
}

func (p *Pam) sha512Sum(data []byte) string {
	bytes := sha512.Sum512(data)
	return hex.EncodeToString(bytes[:])
}

func (p *Pam) GenerateHash(hashType, user, password string, salt bool, strething int) string {
	var m HashMethod
	var h string

	switch hashType {
	case "sha512":
		m = p.sha512Sum
	default:
		m = p.sha256Sum
	}

	if salt {
		h = m([]byte(user))
	}

	h = strings.ToLower(m([]byte(h + password)))

	for i := 0; i < strething; i++ {
		h = strings.ToLower(m([]byte(h)))
	}
	return h
}

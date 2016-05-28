package hash

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

type HashMethod func([]byte) string

func sha256Sum(data []byte) string {
	bytes := sha256.Sum256(data)
	return hex.EncodeToString(bytes[:])
}

func sha512Sum(data []byte) string {
	bytes := sha512.Sum512(data)
	return hex.EncodeToString(bytes[:])
}

func Calculate(hashType string, salt bool, user string, password string, strething int) string {
	var m HashMethod
	var h string

	switch hashType {
	case "sha512":
		m = sha512Sum
	default:
		m = sha256Sum
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

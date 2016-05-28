package main

import (
	"testing"

	"github.com/STNS/libnss_stns/config"
)

func TestAuthOk(t *testing.T) {
	config, _ := config.Load("./fixtures/auth_01.conf")
	if checkPassword(config, "sudo", "example", "test") != PAM_SUCCESS {
		t.Error("auth error auth ok")
	}
}

func TestAuthNg(t *testing.T) {
	config, _ := config.Load("./fixtures/auth_01.conf")

	if checkPassword(config, "sudo", "example", "notmatch") != PAM_AUTH_ERR {
		t.Error("auth error auth ng")
	}
}

func TestSalt(t *testing.T) {
	config, _ := config.Load("./fixtures/auth_02.conf")

	if checkPassword(config, "sudo", "example", "test") != PAM_SUCCESS {
		t.Error("auth error salt")
	}
}

func TestStretching(t *testing.T) {
	config, _ := config.Load("./fixtures/auth_03.conf")

	if checkPassword(config, "sudo", "example", "test") != PAM_SUCCESS {
		t.Error("auth error stretching")
	}
}

func TestException(t *testing.T) {
	config, _ := config.Load("./fixtures/auth_04.conf")

	if checkPassword(config, "sudo", "example", "test") != PAM_AUTHINFO_UNAVAIL {
		t.Error("auth error exeption")
	}
}

func TestHashType(t *testing.T) {
	config, _ := config.Load("./fixtures/auth_05.conf")

	if checkPassword(config, "sudo", "global", "test") != PAM_SUCCESS {
		t.Error("auth error hash1")
	}

	if checkPassword(config, "sudo", "example", "test") != PAM_SUCCESS {
		t.Error("auth error hash2")
	}
}

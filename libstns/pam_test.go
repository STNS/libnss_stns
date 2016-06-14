package libstns

import "testing"

func TestAuthOk(t *testing.T) {
	config, _ := LoadConfig("./pam_fixtures/auth_01.conf")
	pam := NewPam(config, 2, []string{"sudo", "example"})
	if pam.PasswordAuth("example", "test") != PAM_SUCCESS {
		t.Error("auth error auth ok")
	}
}

func TestAuthNg(t *testing.T) {
	config, _ := LoadConfig("./pam_fixtures/auth_01.conf")
	pam := NewPam(config, 2, []string{"sudo", "example"})
	if pam.PasswordAuth("example", "nomatch") != PAM_AUTH_ERR {
		t.Error("auth error auth ng")
	}
}

func TestSalt(t *testing.T) {
	config, _ := LoadConfig("./pam_fixtures/auth_02.conf")

	pam := NewPam(config, 2, []string{"sudo", "example"})
	if pam.PasswordAuth("example", "test") != PAM_SUCCESS {
		t.Error("auth error salt")
	}
}

func TestStretching(t *testing.T) {
	config, _ := LoadConfig("./pam_fixtures/auth_03.conf")
	pam := NewPam(config, 2, []string{"sudo", "example"})
	if pam.PasswordAuth("example", "test") != PAM_SUCCESS {
		t.Error("auth error stretching")
	}
}

func TestException(t *testing.T) {
	config, _ := LoadConfig("./pam_fixtures/auth_04.conf")
	pam := NewPam(config, 2, []string{"sudo", "example"})
	if pam.PasswordAuth("example", "test") != PAM_AUTHINFO_UNAVAIL {
		t.Error("auth error exeption")
	}
}

func TestHashType(t *testing.T) {
	config, _ := LoadConfig("./pam_fixtures/auth_05.conf")
	pam := NewPam(config, 2, []string{"sudo", "example"})
	// use metadata hash type
	if pam.PasswordAuth("global", "test") != PAM_SUCCESS {
		t.Error("auth error hash1")
	}

	// use user hash type
	if pam.PasswordAuth("example", "test") != PAM_SUCCESS {
		t.Error("auth error hash1")
	}
}

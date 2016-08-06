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

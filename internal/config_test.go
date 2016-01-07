package libnss_stns

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	configFile, err := ioutil.TempFile("", "libnss_stns-config-test")
	assertNoError(t, err)

	configContent := "api_end_point=\"is string\""

	_, err = configFile.WriteString(configContent)
	assertNoError(t, err)

	configFile.Close()
	defer os.Remove(configFile.Name())

	config, err := LoadConfig(configFile.Name())
	assertNoError(t, err)
	assert(t, config.ApiEndPoint == "is string", "ng api endpoint")
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

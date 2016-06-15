package libstns

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/STNS/libnss_stns/test"
)

func TestLoadConfig(t *testing.T) {

	configFile, err := ioutil.TempFile("", "libnss_stns-config-test")
	test.AssertNoError(t, err)

	configContent := "api_end_point=[\"is string\", \"is string\"]"

	_, err = configFile.WriteString(configContent)
	test.AssertNoError(t, err)

	configFile.Close()
	defer os.Remove(configFile.Name())

	config, err := LoadConfig(configFile.Name())
	test.AssertNoError(t, err)
	test.Assert(t, config.ApiEndPoint[0] == "is string", "ng api endpoint")
	test.Assert(t, config.ApiEndPoint[1] == "is string", "ng api endpoint")
}

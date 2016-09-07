package libstns

import (
	"testing"

	"github.com/STNS/libnss_stns/test"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./fixtures/config/test_config_001.conf")
	test.AssertNoError(t, err)
	test.Assert(t, config.ApiEndPoint[0] == "http://api01.example.com", "ng api endpoint1")
	test.Assert(t, config.ApiEndPoint[1] == "http://api02.example.com", "ng api endpoint2")
	test.Assert(t, config.RequestHeader["x-api-key"] == "fuga", "ng request header")
	test.Assert(t, config.RequestRetry == 3, "ng request retry")
}

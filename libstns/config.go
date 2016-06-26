package libstns

import (
	"github.com/BurntSushi/toml"
	"github.com/STNS/libnss_stns/settings"
)

type Config struct {
	ApiEndPoint     []string `toml:"api_end_point"`
	RequestTimeOut  int      `toml:"request_timeout"`
	User            string   `toml:"user"`
	Password        string   `toml:"password"`
	SslVerify       bool     `toml:"ssl_verify"`
	WrapperCommand  string   `toml:"wrapper_path"`
	ChainSshWrapper string   `toml:"chain_ssh_wrapper"`
	HttpProxy       string   `toml:"http_proxy"`
}

func LoadConfig(filePath string) (*Config, error) {
	var config Config

	defaultConfig(&config)
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func defaultConfig(config *Config) {
	config.RequestTimeOut = settings.HTTP_TIMEOUT
	config.WrapperCommand = "/usr/local/bin/stns-query-wrapper"
	config.ApiEndPoint = []string{"http://localhost:1104"}
}

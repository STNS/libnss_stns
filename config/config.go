package config

import "github.com/pyama86/toml"

type Config struct {
	ApiEndPoint     []string `toml:"api_end_point"`
	User            string   `toml:"user"`
	Password        string   `toml:"password"`
	SslVerify       bool     `toml:"ssl_verify"`
	ChainSshWrapper string   `toml:"chain_ssh_wrapper"`
}

func Load(filePath string) (*Config, error) {
	var config Config

	defaultConfig(&config)
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}

func defaultConfig(config *Config) {
	config.ApiEndPoint = []string{"http://localhost:1104"}
}

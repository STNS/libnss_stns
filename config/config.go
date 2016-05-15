package config

import "github.com/BurntSushi/toml"

type Config struct {
	ApiEndPoint     []string `toml:"api_end_point"`
	User            string   `toml:"user"`
	Password        string   `toml:"password"`
	SslVerify       bool     `toml:"ssl_verify"`
	WrapperCommand  string   `toml:"wrapper_path"`
	ChainSshWrapper string   `toml:"chain_ssh_wrapper"`
}

func Load(filePath string) (*Config, error) {
	var config Config

	defaultConfig(&config)
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func defaultConfig(config *Config) {
	config.WrapperCommand = "/usr/local/bin/stns-query-wrapper"
	config.ApiEndPoint = []string{"http://localhost:1104"}
}

package config

import "github.com/BurntSushi/toml"

type Config struct {
	ApiEndPoint     []string `toml:"api_end_point"`
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

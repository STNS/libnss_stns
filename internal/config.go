package libnss_stns

import "github.com/BurntSushi/toml"

type Config struct {
	Api_End_Point string
}

func LoadConfig(filePath string) (*Config, error) {
	var config Config

	// is unit test
	if filePath != "/etc/stns/libnss_stns.conf" {
		config.Api_End_Point = filePath
		return &config, nil
	}

	defaultConfig(&config)
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}

func defaultConfig(config *Config) {
	config.Api_End_Point = "http://localhost:1104"
}

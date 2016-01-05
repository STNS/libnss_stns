package libnss_stns

import "github.com/BurntSushi/toml"

const filePath = "/etc/stns/libnss_stns.conf"

type Config struct {
	Api_End_Point string
}

func LoadConfig() *Config {
	var config Config
	defaultConfig(&config)
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func defaultConfig(config *Config) {
	config.Api_End_Point = "http://localhost:1104"
}

package main

import "github.com/BurntSushi/toml"

const filePath = "/etc/libnss_etcd.conf"

type Config struct {
	EtcdUrl       string
	UserEndPoint  string
	GroupEndPoint string
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
	config.EtcdUrl = "http://localhost:2379"
	config.UserEndPoint = "users"
	config.GroupEndPoint = "groups"
}

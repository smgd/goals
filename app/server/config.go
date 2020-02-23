package server

import (
	"goals/app/store"
)


// Server's config
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}


// Server's Config constructor
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Store:	  store.NewConfig()	
	}
}
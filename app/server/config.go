package server

import (
	"goals/app/store"
)

// Server's config
type Config struct {
	BindAddr        string `toml:"bind_addr"`
	TokenSigningKey string `toml:"token_signing_key"`
	LogLevel        string `toml:"log_level"`
	Store           *store.Config
}

// Server's Config constructor
func NewConfig() *Config {
	return &Config{Store: store.NewConfig()}
}

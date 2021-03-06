package main

import (
	"flag"

	"github.com/smgd/goals/app/server"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to server's config file")
}

func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		panic(err)
	}

	if err := server.Start(config); err != nil {
		panic(err)
	}
}

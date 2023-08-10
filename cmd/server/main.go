package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/GlebSr/nuhaiShort/internal/app/server"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config")
}

func main() {
	flag.Parse()
	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Start(config))

}

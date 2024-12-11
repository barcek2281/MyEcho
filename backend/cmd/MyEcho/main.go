package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/barcek2281/MyEcho/internal/app/apiserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "config path")
}
func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.NewAPIserver(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

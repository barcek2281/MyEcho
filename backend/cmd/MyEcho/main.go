package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/barcek2281/MyEcho/internal/app/apiserver"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	local      bool
	configPath string
	logPath    string
)

func init() {
	flag.BoolVar(&local, "local", false, "use for start localhost")
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "config path")
	flag.StringVar(&logPath, "log-path", "log/info.log", "log path")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	if local {
		config.BinAddr = "localhost:8080"
	}
	fmt.Printf("http://%v\n", config.BinAddr)

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/barcek2281/MyEcho/internal/app/apiserver"
)

var (
	configPath string
	logPath    string
)

func init() {
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

	env := apiserver.NewEnv()
	_, err = toml.DecodeFile(configPath, env)
	if err != nil {
		log.Fatal(err)
	}

	env.LogFilePath = logPath

	fmt.Println("http://localhost" + config.BinAddr)

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}

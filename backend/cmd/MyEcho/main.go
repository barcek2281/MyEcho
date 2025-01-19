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

	// fmt.Println(config.DataBaseURL)
	// m, err := migrate.New(
	// 	"./migrations", // Путь к папке с миграциями
	// 	config.DataBaseURL,  // Строка подключения
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Applying migrations...")
	// if err := m.Up(); err != nil {
	// 	if err.Error() == "no change" {
	// 		fmt.Println("No migrations to apply")
	// 	} else {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	fmt.Println("Migrations applied successfully!")
	// }

	fmt.Println("http://localhost" + config.BinAddr)

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}

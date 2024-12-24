package apiserver

//import "github.com/barcek2281/MyEcho/internal/app/storage"

type Config struct {
	BinAddr         string `toml:"bin_addr"`
	DataBaseURL     string `toml:"database_url"`
	CookieKey       string `toml:"cookie_key"`
}

func NewConfig() *Config {
	return &Config{
		BinAddr:     ":8080",
		DataBaseURL: "postgres://postgres:admin@localhost:5432/api",
		CookieKey:   "Cookie",
	}
}

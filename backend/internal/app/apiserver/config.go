package apiserver

// import "github.com/barcek2281/MyEcho/internal/app/storage"
type Config struct {
	BinAddr         string `toml:"bin_addr"`
	DataBaseURL     string `toml:"database_url"`
	CookieKey       string `toml:"cookie_key"`
	EmailTo         string `toml:"email"`
	EmailToPassword string `toml:"email_password"`
	LogLevel        string `toml:"log_level"`
	LogFilePath     string `toml:"log_file"`
}

func NewConfig() *Config {
	return &Config{
		BinAddr:         ":8080",
		DataBaseURL:     "postgres://postgres:admin@localhost:5432/api",
		CookieKey:       "Cookie",
		EmailTo:         "sabdpp17@gmail.com",
		EmailToPassword: "123 456 789 0123",
		LogLevel:        "debug",
		LogFilePath:     "./log/info.log",
	}
}

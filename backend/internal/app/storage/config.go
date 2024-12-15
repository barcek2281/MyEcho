package storage

type Config struct {
	DataBaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		DataBaseURL: "postgres://postgres:admin@localhost:5432/api",
	}
}

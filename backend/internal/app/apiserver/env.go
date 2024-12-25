package apiserver

type Env struct {
	EmailTo         string `toml:"email"`
	EmailToPassword string `toml:"email_password"`
	LogLevel        string `toml:"log_level"`
	LogFilePath		string `toml:"log_file"`
}

func NewEnv() *Env {
	return &Env{
		"example@mai.com",
		"examplePassword",
		"debug",
		"logfilePath",
	}
}

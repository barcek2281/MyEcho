package apiserver

type Env struct {
	EmailTo         string `toml:"email"`
	EmailToPassword string `toml:"email_password"`
}

func NewEnv() *Env {
	return &Env{
		"example@mai.com",
		"examplePassword",
	}
}

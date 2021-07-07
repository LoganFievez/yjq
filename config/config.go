package config

type Config struct {
	Filename string
	Pipemode bool
}

func New() Config {
	config := Config{
		Pipemode: false,
		Filename: "",
	}
	return config
}

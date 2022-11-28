package config

type Config struct {
	Verbosity int
}

func NewConfig() Config {
	return Config{}
}

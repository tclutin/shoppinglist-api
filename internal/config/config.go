package config

type Config struct {
}

func MustLoad() *Config {
	var config Config

	return &config
}

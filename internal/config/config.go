package config

import "time"

type Config struct {
	HTTPServer HTTPServer
	Postgres   Postgres
	JWT        JWT
}

type HTTPServer struct {
	Host string
	Port string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type JWT struct {
	Secret string
	Expire time.Duration
}

func MustLoad() *Config {
	var config Config

	return &config
}

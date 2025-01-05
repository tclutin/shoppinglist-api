package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

const (
	dev  string = "dev"
	prod string = "prod"
)

type Config struct {
	Env        string `env:"env"`
	HTTPServer HTTPServer
	Postgres   Postgres
	JWT        JWT
}

type HTTPServer struct {
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DATABASE"`
}

type JWT struct {
	Secret        string        `env:"JWT_SECRET"`
	AccessExpire  time.Duration `env:"JWT_ACCESS_EXPIRE"`
	RefreshExpire time.Duration `env:"JWT_REFRESH_EXPIRE"`
}

func MustLoad() *Config {
	var config Config

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalln("Error reading env", err)
	}

	return &config
}

func (c Config) IsProd() bool {
	return c.Env == prod
}

func (c Config) IsDev() bool {
	return c.Env == dev
}

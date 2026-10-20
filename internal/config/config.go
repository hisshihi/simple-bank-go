package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBHost string `env:"DB_HOST"     envDefault:"localhost"`
	DBPort int    `env:"DB_PORT"     envDefault:"5432"`
	DBUser string `env:"DB_USER,required"`
	DBPass string `env:"DB_PASS,required"`
	DBName string `env:"DB_NAME,required"`

	//	Logger
	LoggerLevel  string `env:"LOGGER_LEVEL"`
	LoggerPretty bool   `env:"LOGGER_PRETTY"`

	HTTPAddr string `env:"HTTP_ADDR"`
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
	)
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	return &cfg, nil
}

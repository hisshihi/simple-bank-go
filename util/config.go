// Package util содержит утилиты и вспомогательные функции
package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`

	DBSimpleBankUser     string `mapstructure:"DB_SIMPLE_BANK_USER"`
	DBSimpleBankPassword string `mapstructure:"DB_SIMPLE_BANK_PASSWORD"`
	DBSimpleBankDatabase string `mapstructure:"DB_SIMPLE_BANK_DATABASE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return
}

package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DbName   string
		SslMode  string
	}
}

func Load() *Config {
	viper.
}
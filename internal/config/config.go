package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	wd, err := os.Getwd()
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Error getting executable path, %s", err))
	}

	projectRoot := findProjectRoot(wd)
	configPath := filepath.Join(projectRoot, "config")

	viper.SetConfigName("config.dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Fatal(fmt.Sprintf("Error reading config file, %s", err))
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		zap.L().Fatal(fmt.Sprintf("Unable to decode into struct, %v", err))
	}
	return &config
}

func findProjectRoot(startPath string) string {
	dir := startPath
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			zap.L().Fatal("Could not find project root with go.mod")
		}
		dir = parentDir
	}
}

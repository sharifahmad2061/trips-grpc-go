package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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
		log.Fatalf("Error getting executable path, %s", err)
	}

	projectRoot := findProjectRoot(wd)
	configPath := filepath.Join(projectRoot, "config")

	viper.SetConfigName("config.dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
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
			log.Fatal("Could not find project root with go.mod")
		}
		dir = parentDir
	}
}

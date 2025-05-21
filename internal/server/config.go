package server

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

type Config struct {
	ServAddr    string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func NewConfig() Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse env vars: %v", err)
	}
	flag.StringVar(&config.ServAddr, "a", config.ServAddr, "Server address (default from env: SERVER_ADDRESS)")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "Base URL (default from env: BASE_URL)")
	flag.StringVar(&config.FileStorage, "f", config.FileStorage, "File Storage (default from env: FILE_STORAGE_PATH)")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "Database DSN (default from env: DATABASE_DSN)")

	flag.Parse()
	if config.ServAddr == "" {
		config.ServAddr = "localhost:8080"
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}
	return config
}

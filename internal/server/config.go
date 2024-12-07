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
}

func NewConfig() *Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("Failed to parse env vars: %v", err)
	}

	if config.ServAddr == "" {
		flag.StringVar(&config.ServAddr, "a", "localhost:8080", "serv address")
	}
	if config.BaseURL == "" {
		flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base URL")
	}
	if config.FileStorage == "" {
		flag.StringVar(&config.FileStorage, "f", "", "File Storage")

	}
	flag.Parse()
	return &config
}

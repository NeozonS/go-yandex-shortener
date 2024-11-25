package server

type Config struct {
	ServAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

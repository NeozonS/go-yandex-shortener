package server

type Config struct {
	ServAddr string `env:"SERV_ADDR" envDefault:":8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost"`
}

package main

import (
	"flag"
	"github.com/NeozonS/go-shortener-ya.git/internal/handlers"
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/file"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/mapbd"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	config := server.Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("Failed to parse env vars: %v", err)
	}
	pars_argg(&config)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	repositories := choiseStorage(config.FileStorage)
	handler := handlers.NewHandlers(repositories, config)
	r.Route("/", func(r chi.Router) {
		r.Post("/api/shorten", handler.PostAPI)
		r.Post("/", handler.PostHandler)
		r.Get("/{id}", handler.GetHandler)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		})
	})
	log.Println("Server started at " + config.ServAddr)
	http.ListenAndServe(config.ServAddr, r)
}

func choiseStorage(storage string) storage.Repositories {
	if storage == "" {
		return mapbd.New()
	}
	return file.NewFileStorage(storage)
}
func pars_argg(config *server.Config) *server.Config {
	if config.ServAddr == "localhost:8080" {
		flag.StringVar(&config.ServAddr, "a", config.ServAddr, "serv address")
	}
	if config.BaseURL == "http://localhost:8080" {
		flag.StringVar(&config.BaseURL, "b", config.BaseURL, "Base URL")
	}
	if config.FileStorage == "" {
		flag.StringVar(&config.FileStorage, "f", config.FileStorage, "File Storage")

	}
	flag.Parse()
	return config
}

package main

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/handlers"
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	repositories := mapbd.New()
	handler := handlers.NewHandlers(repositories, config)
	r.Route("/", func(r chi.Router) {
		r.Post("/api/shorten", handler.PostAPI)
		r.Post("/", handler.PostHandler)
		r.Get("/{id}", handler.GetHandler)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		})
	})
	log.Printf("Using config: %+v", config)
	log.Println("Server started at " + config.BaseURL + config.ServAddr)
	http.ListenAndServe(":"+config.ServAddr, r)
}

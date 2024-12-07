package main

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/handlers"
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/file"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/mapbd"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	config := server.NewConfig()
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

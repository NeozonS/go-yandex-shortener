package main

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/handlers"
	"github.com/NeozonS/go-shortener-ya.git/internal/middleware"
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/file"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/mapbd"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	config := server.NewConfig()

	repositories, err := choiseStorage(config.FileStorage)
	if err != nil {
		log.Fatal(err)
	}
	handler := handlers.NewHandlers(repositories, config)

	r := chi.NewRouter()
	r.Use(middleware.CookieMiddleware)
	r.Use(middleware.GzipRequestMiddleware)
	r.Use(middleware.GzipResponseMiddleware)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handler.PostAPI)
		r.Get("/user/urls", handler.GetAPIAllURLHandler)
	})

	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	log.Println("Server started at " + config.ServAddr)
	http.ListenAndServe(config.ServAddr, r)
}

func choiseStorage(storage string) (storage.Repository, error) {
	if storage == "" {
		return mapbd.New()
	}
	if storage == "POSTGRES" {
		dsn := "postgres://postgres:123456l@localhost:5432/shortener_db?sslmode=disable"
		return postgres.NewPostgresDB(dsn)
	}
	return file.NewFileStorage(storage)
}

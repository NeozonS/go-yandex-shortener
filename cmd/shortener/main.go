package main

import (
	"context"
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

	repositories, err := choiseStorage(config)
	if err != nil {
		log.Fatal(err)
	}
	if pgStore, ok := repositories.(*postgres.PostgresDB); ok {
		if err := pgStore.CreateTable(context.Background()); err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
	}

	handler := handlers.NewHandlers(repositories, config)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)
	r.Use(middleware.GzipRequestMiddleware)
	r.Use(middleware.GzipResponseMiddleware)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handler.PostAPI)
		r.Get("/user/urls", handler.GetAPIAllURLHandler)
		r.Get("/shorten/batch", handler.BatchHandler)
	})

	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	r.Get("/ping", handler.PingHandler)
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	log.Println("Server started at " + config.ServAddr)
	http.ListenAndServe(config.ServAddr, r)
}

func choiseStorage(storage server.Config) (storage.Repository, error) {
	if storage.DatabaseDSN != "" {
		store, err := postgres.NewPostgresDB(storage.DatabaseDSN)
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
		}
		log.Println("Using PostgreSQL storage")
		return store, nil
	} else if storage.FileStorage != "" {
		store, err := file.NewFileStorage(storage.FileStorage)
		if err != nil {
			log.Fatalf("Failed to initialize file storage: %v", err)
		}
		log.Println("Using file storage")
		return store, nil
	} else {
		store, _ := mapbd.New()
		log.Println("Using mapbd storage")
		return store, nil
	}
	//	dsn := "postgres://postgres:123456l@localhost:5432/shortener_db?sslmode=disable"

}

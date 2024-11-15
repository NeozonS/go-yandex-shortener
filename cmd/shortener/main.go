package main

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/handlers"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/mapbd"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	repositories := mapbd.New()
	handler := handlers.NewHandlers(repositories)
	r.Route("/", func(r chi.Router) {
		r.Post("/api/shorten", handler.PostAPI)
		r.Post("/", handler.PostHandler)
		r.Get("/{id}", handler.GetHandler)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		})
	})

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

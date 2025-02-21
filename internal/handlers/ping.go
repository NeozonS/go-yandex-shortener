package handlers

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/postgres"
	"net/http"
)

type DB struct {
	DB *postgres.PostgresDB
}

func (u *Handlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := u.repo.Ping(r.Context())
	if err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

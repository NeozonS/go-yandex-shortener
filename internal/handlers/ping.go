package handlers

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/postgres"
	"net/http"
)

type HandlersTest struct {
	DB *postgres.PostgresDB
}

func (h *HandlersTest) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.DB.DB.Ping()
	if err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

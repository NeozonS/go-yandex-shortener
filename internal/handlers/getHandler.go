package handlers

import (
	"database/sql"
	"errors"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (u *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "id")
	if token == "" {
		http.Error(w, "Запрашиваемая страница не найдена", 400)
		return
	}
	var (
		originalURL string
		err         error
		isDeleted   bool
	)
	if dr, ok := u.repo.(postgres.DeletableRepository); ok {
		originalURL, isDeleted, err = dr.GetURLWithDeleted(r.Context(), token)
	} else {
		originalURL, err = u.repo.GetURL(r.Context(), token)
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	case isDeleted:
		w.WriteHeader(http.StatusGone)
		return
	}
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(307)

}

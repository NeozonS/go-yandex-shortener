package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (u *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlP := chi.URLParam(r, "id")
	originalURL, err := u.repo.GetURL(u.config.BaseURL + "/" + urlP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if urlP == "" {
		http.Error(w, "Запрашиваемая страница не найдена", 400)
		return
	}
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(307)

}

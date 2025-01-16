package handlers

import (
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"io"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func (u *Handlers) PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	originURL := string(b)
	if originURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(string(b), "http://") && !strings.HasPrefix(string(b), "https://") {
		originURL = "http://" + originURL
	}
	shortURL := u.config.BaseURL + "/" + utils.GenerateShortURL()
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	err = u.repo.UpdateURL(userID, shortURL, originURL)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.WriteHeader(201)
	fmt.Fprint(w, shortURL)
}

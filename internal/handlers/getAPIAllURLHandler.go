package handlers

import (
	"encoding/json"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"net/http"
)

func (u *Handlers) GetAPIAllURLHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}
	allURL, err := u.repo.GetAllURL(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
	}

	for i, link := range allURL {
		allURL[i].ShortURL = utils.FullURL(u.config.BaseURL, link.ShortURL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(allURL)
}

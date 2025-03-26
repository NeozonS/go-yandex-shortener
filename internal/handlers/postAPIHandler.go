package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"net/http"
	"strings"
)

func (u *Handlers) PostAPI(w http.ResponseWriter, r *http.Request) {
	b := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var jsonurl APIJson
	err := b.Decode(&jsonurl)
	if !strings.HasPrefix(jsonurl.URL, "http://") && !strings.HasPrefix(jsonurl.URL, "https://") {
		jsonurl.URL = "http://" + jsonurl.URL
	}
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	userID, ok := utils.GetUserID(r.Context())
	token := utils.GenerateShortURL(jsonurl.URL, userID)
	if !ok || userID == "" {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}
	err = u.repo.UpdateURL(r.Context(), userID, token, jsonurl.URL)
	var conflictErr models.ErrURLConflict
	if err != nil {
		if errors.As(err, &conflictErr) {
			w.WriteHeader(409)
			json.NewEncoder(w).Encode(APIJson{Result: utils.FullURL(u.config.BaseURL, conflictErr.ExistingURL)})
			return
		}
	}
	result := APIJson{Result: utils.FullURL(u.config.BaseURL, token)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
}

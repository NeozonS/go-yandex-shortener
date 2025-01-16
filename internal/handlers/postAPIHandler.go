package handlers

import (
	"encoding/json"
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
	shortURL := u.config.BaseURL + "/" + utils.GenerateShortURL()
	userID, ok := utils.GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}
	err = u.repo.UpdateURL(userID, shortURL, jsonurl.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	result := APIJson{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
}

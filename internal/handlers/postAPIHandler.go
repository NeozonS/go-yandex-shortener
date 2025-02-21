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
	userID, ok := utils.GetUserID(r.Context())
	token := utils.GenerateShortURL(jsonurl.URL, userID)
	if !ok || userID == "" {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}
	err = u.repo.UpdateURL(userID, token, jsonurl.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	result := APIJson{Result: utils.FullURL(u.config.BaseURL, token)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
}

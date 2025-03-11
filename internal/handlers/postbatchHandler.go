package handlers

import (
	"encoding/json"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"net/http"
)

func (u *Handlers) PostBatchHandler(w http.ResponseWriter, r *http.Request) {
	var BRequest []BatchRequest
	var BResponse []BatchResponse
	URLs := make(map[string]string)
	ctx := r.Context()
	b := json.NewDecoder(r.Body)
	err := b.Decode(&BRequest)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	userID, ok := utils.GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}
	for _, s := range BRequest {
		token := utils.GenerateShortURL(s.OriginalURL, userID)
		URLs[token] = s.OriginalURL
		BResponse = append(BResponse, BatchResponse{CorrelationID: s.CorrelationID, ShortURL: utils.FullURL(u.config.BaseURL, token)})
	}
	u.repo.BatchUpdateURL(ctx, userID, URLs)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(BResponse)
}

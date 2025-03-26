package handlers

import (
	"encoding/json"
	"errors"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
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
	err = u.repo.BatchUpdateURL(ctx, userID, URLs)
	switch {
	case errors.Is(err, models.ErrURLConflict):
		w.WriteHeader(http.StatusConflict)
	case err != nil:
		http.Error(w, err.Error(), 400)
	default:
		w.WriteHeader(201)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BResponse)
}

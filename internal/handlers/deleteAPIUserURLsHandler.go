package handlers

import (
	"encoding/json"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"net/http"
)

func (u *Handlers) DeleteAPIUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}
	var ids []string
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, id := range ids {
		u.worker.EnqueueURLForDeleteion(userID, id)
	}
	w.WriteHeader(http.StatusAccepted)

}

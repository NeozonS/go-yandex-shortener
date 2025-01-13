package handlers

import (
	"encoding/json"
	"net/http"
)

func (u *Handlers) GetAPIAllURLHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	allURL, err := u.repo.GetAllURL(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(allURL)

}

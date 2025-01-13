package utils

import (
	"errors"
	"net/http"
)

func SetCookie(w http.ResponseWriter, userID string) {
	signature, err := Encrypt(userID)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "userID",
		Value:  signature,
		Path:   "/",
		MaxAge: 3600, // Кука будет действительна 1 час
	})

}

func GetUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("userID")
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}
	userID, err := Decrypt(cookie.Value)
	if err != nil {
		return "", err
	}
	return userID, nil
}

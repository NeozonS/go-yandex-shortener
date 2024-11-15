package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Handlers struct {
	repo storage.Repositories
}
type APIJson struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}

func NewHandlers(repo storage.Repositories) *Handlers {
	return &Handlers{repo}
}
func (u *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlP := chi.URLParam(r, "id")
	originalURL, err := u.repo.GetURL(urlP)
	if urlP == "" || err != nil {
		http.Error(w, "Запрашиваемая страница не найдена", 400)
		return
	}
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(307)

}

func (u *Handlers) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод для ПостЗапроса, не правильно указан метод.", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortKey := u.generateShortURL()
	err = u.repo.UpdateURL(string(b), shortKey)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprint(w, "http://localhost:8080/"+shortKey)
}

func (u *Handlers) PostAPI(w http.ResponseWriter, r *http.Request) {
	b := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var Url APIJson
	err := b.Decode(&Url)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	shortKey := u.generateShortURL()
	err = u.repo.UpdateURL(Url.URL, shortKey)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	result := APIJson{Result: "http://localhost:8080/" + shortKey}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(result)
}

func (u *Handlers) generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	seed := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seed)
	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = charset[randGen.Intn(len(charset))]
	}
	return string(shortKey)
}

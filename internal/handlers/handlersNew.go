package handlers

import (
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage"
)

type Handlers struct {
	repo   storage.Repository
	config server.Config
}
type APIJson struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}
type BatchRequest struct {
	CorrelationId string `json:"correlation_id,"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func NewHandlers(repo storage.Repository, config server.Config) *Handlers {
	return &Handlers{repo, config}
}

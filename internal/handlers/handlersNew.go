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

func NewHandlers(repo storage.Repository, config server.Config) *Handlers {
	return &Handlers{repo, config}
}

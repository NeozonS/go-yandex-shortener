package storage

import "github.com/NeozonS/go-shortener-ya.git/internal/storage/models"

type Repository interface {
	GetURL(shortURL string) (string, error)
	UpdateURL(userID, shortURL, originalURL string) error
	GetAllURL(userID string) ([]models.LinkPair, error)
}

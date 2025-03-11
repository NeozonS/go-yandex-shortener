package storage

import (
	"context"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
)

type Repository interface {
	GetURL(ctx context.Context, shortURL string) (string, error)
	UpdateURL(ctx context.Context, userID, shortURL, originalURL string) error
	BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error
	GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error)
	Ping(ctx context.Context) error
}

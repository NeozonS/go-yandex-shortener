package mapbd

import (
	"context"
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
)

type MapBD struct {
	Urls map[string]map[string]string
}

func (m *MapBD) GetURL(ctx context.Context, shortURL string) (string, error) {
	for _, u := range m.Urls {
		if originalURL, ok := u[shortURL]; ok {
			return originalURL, nil
		}
	}
	return "", fmt.Errorf("short url not found for %s", shortURL)
}

func (m *MapBD) GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error) {
	u, ok := m.Urls[userID]
	if !ok {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	userLink := make([]models.LinkPair, 0, len(u))
	for shortURL, originalURL := range u {
		userLink = append(userLink, models.LinkPair{ShortURL: shortURL, LongURL: originalURL})
	}
	if len(userLink) == 0 {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	return userLink, nil
}

func (m *MapBD) UpdateURL(ctx context.Context, userID, shortURL, originalURL string) error {
	if _, ok := m.Urls[userID]; !ok {
		m.Urls[userID] = make(map[string]string)
	}
	m.Urls[userID][shortURL] = originalURL
	return nil
}

func (m *MapBD) BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error {
	if _, ok := m.Urls[userID]; !ok {
		m.Urls[userID] = URLs
	}
	m.Urls[userID] = URLs
	return nil
}
func New() (*MapBD, error) {
	return &MapBD{Urls: make(map[string]map[string]string)}, nil
}
func (m *MapBD) Ping(ctx context.Context) error {
	return nil
}

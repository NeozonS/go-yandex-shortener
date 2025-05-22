package mapbd

import (
	"context"
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	"sync"
)

type URLData struct {
	URL     string
	Deleted bool
}

type MapBD struct {
	mu   sync.RWMutex
	Urls map[string]map[string]URLData
}

func (m *MapBD) GetURL(ctx context.Context, shortURL string) (string, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.Urls {
		if originalURL, ok := u[shortURL]; ok {
			if originalURL.Deleted {
				return "", true, nil
			}
			return originalURL.URL, false, nil
		}
	}
	return "", false, fmt.Errorf("short url not found for %s", shortURL)
}

func (m *MapBD) GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error) {
	u, ok := m.Urls[userID]
	if !ok {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	userLink := make([]models.LinkPair, 0, len(u))
	for shortURL, URLData := range u {
		userLink = append(userLink, models.LinkPair{ShortURL: shortURL, LongURL: URLData.URL})

	}
	if len(userLink) == 0 {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	return userLink, nil
}

func (m *MapBD) UpdateURL(ctx context.Context, userID, shortURL, originalURL string) error {
	if _, ok := m.Urls[userID]; !ok {
		m.Urls[userID] = make(map[string]URLData)
	}
	m.Urls[userID][shortURL] = URLData{
		URL:     originalURL,
		Deleted: false,
	}
	return nil
}

func (m *MapBD) BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error {
	for shortURL, originalURL := range URLs {
		m.Urls[userID][shortURL] = URLData{
			URL:     originalURL,
			Deleted: false,
		}
	}
	return nil
}
func New() (*MapBD, error) {
	return &MapBD{Urls: make(map[string]map[string]URLData)}, nil
}
func (m *MapBD) Ping(ctx context.Context) error {
	return nil
}

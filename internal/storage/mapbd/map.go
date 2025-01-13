package mapbd

import (
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
)

type MapBD struct {
	Urls map[string]map[string]string
}

func (m *MapBD) GetURL(shortURL string) (string, error) {
	for _, u := range m.Urls {
		if originalURL, ok := u[shortURL]; ok {
			return originalURL, nil
		}
	}
	return "", fmt.Errorf("short url not found for %s", shortURL)
}
func (m *MapBD) GetAllURL(userID string) ([]models.LinkPair, error) {
	u, ok := m.Urls[userID]
	if !ok {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	userLink := make([]models.LinkPair, 0, len(u))
	for shortURL, originalURL := range u {
		userLink = append(userLink, models.LinkPair{shortURL, originalURL})
	}
	if len(userLink) == 0 {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	return userLink, nil
}

func (m *MapBD) UpdateURL(userID, shortURL, originalURL string) error {
	if _, ok := m.Urls[userID]; !ok {
		m.Urls[userID] = make(map[string]string)
	}
	m.Urls[userID][shortURL] = originalURL
	return nil
}
func New() *MapBD {
	return &MapBD{Urls: make(map[string]map[string]string)}
}

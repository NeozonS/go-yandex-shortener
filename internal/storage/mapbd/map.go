package mapbd

import "fmt"

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
func (m *MapBD) GetAllURL(userID string) (map[string]string, error) {
	u, ok := m.Urls[userID]
	if !ok {
		return nil, fmt.Errorf("user not found for %s", userID)
	}
	return u, nil
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

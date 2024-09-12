package mapbd

import "fmt"

type MapBD struct {
	Urls map[string]string
}

func (m *MapBD) GetURL(id string) (string, error) {
	u, ok := m.Urls[id]
	if !ok {
		return "", fmt.Errorf("url not found for %s", id)
	}
	return u, nil
}
func (m *MapBD) UpdateURL(url, id string) error {
	if url == "" {
		return fmt.Errorf("url not recognized")
	}
	if id == "" {
		return fmt.Errorf("id not created")
	}
	m.Urls[id] = url
	return nil
}
func New() *MapBD {
	return &MapBD{Urls: make(map[string]string)}
}

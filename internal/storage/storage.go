package storage

type Repositories interface {
	GetURL(shortURL string) (string, error)
	UpdateURL(userID, shortURL, originalURL string) error
	GetAllURLs(userID string) (map[string]string, error)
}

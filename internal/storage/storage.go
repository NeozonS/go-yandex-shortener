package storage

type Repositories interface {
	GetURL(url string) (string, error)
	UpdateURL(url, id string) error
}

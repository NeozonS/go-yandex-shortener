package models

type LinkPair struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

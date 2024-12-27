package file

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Storage struct {
	file *os.File
}
type UserUrl struct {
	UserID string     `json:"user_id"`
	Links  []LinkPair `json:"links"`
}
type LinkPair struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

func (m *Storage) GetURL(shortURL string) (string, error) {
	file, err := os.Open(m.file.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for {
		var pair LinkPair
		err := decoder.Decode(&pair)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if pair.ShortURL == shortURL {
			return pair.LongURL, nil
		}
	}
	return "", errors.New("url not found")
}
func (m *Storage) GetAllURL(userID string) ([]LinkPair, error) {
	file, err := os.Open(m.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var pair UserUrl
		err := decoder.Decode(&pair)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if pair.UserID == userID {
			return pair.Links, nil
		}
	}
	return nil, errors.New("user ID not found")
}

func (m *Storage) UpdateURL(userID, shortURL, originalURL string) error {
	file, err := os.OpenFile(m.file.Name(), os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	pair := LinkPair{shortURL, originalURL}
	user := UserUrl{UserID: userID, Links: []LinkPair{pair}}

	encoder := json.NewEncoder(file)
	return encoder.Encode(&user)
}

func NewFileStorage(filename string) *Storage {
	//file, err := os.Create(filename)
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil
	}
	defer file.Close()
	return &Storage{file: file}
}

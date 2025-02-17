package file

import (
	"encoding/json"
	"errors"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	"io"
	"os"
)

type Storage struct {
	file *os.File
}
type UserURL struct {
	UserID string          `json:"user_id"`
	Links  models.LinkPair `json:"links"`
}

func (m *Storage) GetURL(shortURL string) (string, error) {
	file, err := os.Open(m.file.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var user UserURL
		if err := decoder.Decode(&user); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if user.Links.ShortURL == shortURL {
			return user.Links.LongURL, nil
		}
	}
	return "", errors.New("url not found")
}
func (m *Storage) GetAllURL(userID string) ([]models.LinkPair, error) {
	file, err := os.Open(m.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var result []models.LinkPair
	for {
		var pair UserURL
		err := decoder.Decode(&pair)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if pair.UserID == userID {
			result = append(result, pair.Links)
		}
	}
	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New("user ID not found")
}

func (m *Storage) UpdateURL(userID, shortURL, originalURL string) error {
	file, err := os.OpenFile(m.file.Name(), os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	user := UserURL{UserID: userID, Links: models.LinkPair{ShortURL: shortURL, LongURL: originalURL}}

	encoder := json.NewEncoder(file)
	return encoder.Encode(&user)
}

func NewFileStorage(filename string) (*Storage, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &Storage{file: file}, nil
}

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
	UserID string            `json:"user_id"`
	Links  []models.LinkPair `json:"links"`
}

func (m *Storage) GetURL(shortURL string) (string, error) {
	_, err := m.file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	decoder := json.NewDecoder(m.file)
	for {
		var pair models.LinkPair
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
func (m *Storage) GetAllURL(userID string) ([]models.LinkPair, error) {
	file, err := m.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

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
			result = append(result, pair.Links...)
		}
	}
	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New("user ID not found")
}

func (m *Storage) UpdateURL(userID, shortURL, originalURL string) error {
	file, err := m.file.Seek(0, 2)
	if err != nil {
		return err
	}

	pair := models.LinkPair{ShortURL: shortURL, LongURL: originalURL}
	user := UserURL{UserID: userID, Links: []models.LinkPair{pair}}

	encoder := json.NewEncoder(file)
	return encoder.Encode(&user)
}

func NewFileStorage(filename string) *Storage {
	//file, err := os.Create(filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}
	defer file.Close()
	return &Storage{file: file}
}

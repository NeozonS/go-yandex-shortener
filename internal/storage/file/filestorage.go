package file

import (
	"context"
	"encoding/json"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	"io"
	"os"
	"sync"
)

type Storage struct {
	filename string
	mu       sync.RWMutex
}
type UserURL struct {
	UserID string          `json:"user_id"`
	Links  models.LinkPair `json:"links"`
}

func (m *Storage) GetURL(ctx context.Context, shortURL string) (string, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	file, err := os.Open(m.filename)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var user UserURL
		if err := decoder.Decode(&user); err != nil {
			if err == io.EOF {
				break
			}
			return "", false, err
		}
		if user.Links.ShortURL == shortURL {
			return user.Links.LongURL, user.Links.Deleted, nil
		}
	}
	return "", false, os.ErrNotExist
}
func (m *Storage) GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	file, err := os.Open(m.filename)
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
	if len(result) == 0 {
		return nil, os.ErrNotExist
	}
	return result, nil
}

func (m *Storage) UpdateURL(ctx context.Context, userID, shortURL, originalURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.OpenFile(m.filename, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	user := UserURL{UserID: userID, Links: models.LinkPair{ShortURL: shortURL, LongURL: originalURL, Deleted: false}}

	encoder := json.NewEncoder(file)
	return encoder.Encode(&user)
}
func (m *Storage) BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.OpenFile(m.filename, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	user := UserURL{}
	encoder := json.NewEncoder(file)
	for t, o := range URLs {
		user = UserURL{UserID: userID, Links: models.LinkPair{ShortURL: t, LongURL: o, Deleted: false}}
		if err := encoder.Encode(&user); err != nil {
			return err
		}
	}

	return nil
}

func (m *Storage) BatchDeleteURL(ctx context.Context, userID string, shortURL []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	file, err := os.Open(m.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var allEntries []UserURL
	decode := json.NewDecoder(file)
	for {
		var entry UserURL
		if err := decode.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		allEntries = append(allEntries, entry)
	}

	tokenSet := make(map[string]struct{})
	for _, t := range shortURL {
		tokenSet[t] = struct{}{}
	}
	for i := range allEntries {
		if allEntries[i].UserID == userID {
			if _, ok := tokenSet[allEntries[i].Links.ShortURL]; ok {
				allEntries[i].Links.Deleted = true
			}
		}
	}
	tempFile, err := os.CreateTemp("", "urls_temp_*.txt")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	encoder := json.NewEncoder(tempFile)
	for _, entry := range allEntries {
		if err := encoder.Encode(entry); err != nil {
			return err
		}
	}
	if err := os.Rename(tempFile.Name(), m.filename); err != nil {
		return err
	}
	return nil
}
func NewFileStorage(filename string) (*Storage, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	return &Storage{filename: filename}, nil
}

func (m *Storage) Ping(ctx context.Context) error {
	return nil
}

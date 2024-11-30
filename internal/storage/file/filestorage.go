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
type LinkPair struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func (m *Storage) GetURL(id string) (string, error) {
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

		if pair.ShortURL == id {
			return pair.LongURL, nil
		}
	}
	return "", errors.New("ссылка не найдена")
}
func (m *Storage) UpdateURL(url, id string) error {
	file, err := os.OpenFile(m.file.Name(), os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	pair := LinkPair{id, url}
	encoder := json.NewEncoder(file)
	return encoder.Encode(&pair)
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

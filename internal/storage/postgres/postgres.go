package postgres

import (
	"database/sql"
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func (p *PostgresDB) GetURL(shortURL string) (string, error) {
	// Временная заглушка
	return "", fmt.Errorf("GetURL не реализован")
}

// Заглушка метода для получения всех ссылок пользователя
func (p *PostgresDB) GetAllURL(userID string) ([]models.LinkPair, error) {
	// Временная заглушка
	return nil, fmt.Errorf("GetAllURL не реализован")
}

// Заглушка метода для сохранения новой короткой ссылки
func (p *PostgresDB) UpdateURL(userID, shortURL, originalURL string) error {
	// Временная заглушка
	return fmt.Errorf("UpdateURL не реализован")
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database connection: %s", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging database connection: %s", err)
	}

	return &PostgresDB{DB: db}, nil
}

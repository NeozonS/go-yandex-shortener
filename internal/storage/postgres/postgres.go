package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db *sql.DB
}

func (p *PostgresDB) GetURL(shortURL string) (string, error) {
	query := `
	SELECT original_url
	FROM short_urls
	WHERE token = $1`

	var originalURL string
	err := p.db.QueryRow(query, shortURL).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	return originalURL, nil
}

func (p *PostgresDB) GetAllURL(userID string) ([]models.LinkPair, error) {
	query := `
	SELECT token, original_url
	FROM short_urls
	WHERE user_id = $1
	ORDER BY created_at DESC
`
	rows, err := p.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all URLs: %w", err)
	}
	defer rows.Close()
	var urls []models.LinkPair
	for rows.Next() {
		var url models.LinkPair
		if err := rows.Scan(&url.ShortURL, &url.LongURL); err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return urls, nil
}

func (p *PostgresDB) UpdateURL(userID, shortURL, originalURL string) error {
	if err := p.CreateUser(userID); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	query := `
	INSERT INTO short_urls (token, original_url, user_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (token) DO UPDATE
	SET original_url =EXCLUDED.original_url
`
	_, err := p.db.Exec(query, shortURL, originalURL, userID)
	return err
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database connection: %s", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging database connection: %s", err)
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) CreateTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW()
    );
	CREATE TABLE IF NOT EXISTS short_urls(
	    token CHAR(8) PRIMARY KEY,
	    original_url TEXT NOT NULL,
	    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	    clicks BIGINT DEFAULT 0,
	    created_at TIMESTAMPTZ DEFAULT NOW(),
	    expires_at TIMESTAMPTZ
	);
	
	CREATE INDEX IF NOT EXISTS idx_user_created ON short_urls (user_id, created_at);
`
	_, err := p.db.ExecContext(ctx, query)
	return err
}
func (p *PostgresDB) CreateUser(userid string) error {
	query := `INSERT INTO users (id) VALUES ($1) ON CONFLICT DO NOTHING;`
	_, err := p.db.Exec(query, userid)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}
	return nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

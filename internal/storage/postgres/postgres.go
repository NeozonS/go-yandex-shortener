package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db *sql.DB
}

func (p *PostgresDB) GetURL(ctx context.Context, shortURL string) (string, error) {
	query := `
	SELECT original_url
	FROM short_urls
	WHERE token = $1`

	var originalURL string
	err := p.db.QueryRowContext(ctx, query, shortURL).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	return originalURL, nil
}

func (p *PostgresDB) GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error) {
	query := `
	SELECT token, original_url
	FROM short_urls
	WHERE user_id = $1
	ORDER BY created_at DESC
`
	rows, err := p.db.QueryContext(ctx, query, userID)
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

func (p *PostgresDB) UpdateURL(ctx context.Context, userID, shortURL, originalURL string) error {
	if err := p.CreateUser(userID); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	query := `
	INSERT INTO short_urls (token, original_url, user_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (original_url) DO UPDATE 
	SET original_url = EXCLUDED.original_url
	RETURNING token
`
	var existingToken string
	err := p.db.QueryRowContext(ctx, query, shortURL, originalURL, userID).Scan(&existingToken)
	if err != nil {
		return fmt.Errorf("failed to update URL: %w", err)
	}
	if existingToken != shortURL {
		return models.ErrURLConflict{ExistingURL: existingToken}
	}
	return err
}
func (p *PostgresDB) BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error {
	if err := p.CreateUser(userID); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO short_urls (token, original_url,user_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (token) DO UPDATE
	SET original_url = EXCLUDED.original_url
`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for shortURL, originalURL := range URLs {
		if _, err := stmt.ExecContext(ctx, shortURL, originalURL, userID); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return models.ErrURLConflict{ExistingURL: originalURL}
				}
			}
			return fmt.Errorf("failed to insert URL: %w", err)
		}
	}
	return tx.Commit()
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %s", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database connection: %s", err)
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
	    original_url TEXT NOT NULL UNIQUE,
	    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	    clicks BIGINT DEFAULT 0,
	    created_at TIMESTAMPTZ DEFAULT NOW(),
	    expires_at TIMESTAMPTZ
	);
`
	_, err := p.db.ExecContext(ctx, query)
	return err
}
func (p *PostgresDB) CreateUser(userid string) error {
	query := `INSERT INTO users (id) VALUES ($1) ON CONFLICT DO NOTHING;`
	_, err := p.db.Exec(query, userid)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

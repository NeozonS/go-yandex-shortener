CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE short_urls (
                            token CHAR(8) PRIMARY KEY,
                            original_url TEXT NOT NULL,
                            user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                            clicks BIGINT DEFAULT 0,
                            created_at TIMESTAMPTZ DEFAULT NOW(),
                            expires_at TIMESTAMPTZ
);

CREATE INDEX idx_user_created ON short_urls (user_id, created_at);
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE short_urls (
                                         token CHAR(8) PRIMARY KEY,
                                         original_url TEXT NOT NULL UNIQUE,
                                         user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                         clicks BIGINT DEFAULT 0,
                                         is_delete BOOLEAN,
                                         created_at TIMESTAMPTZ DEFAULT NOW(),
                                         expires_at TIMESTAMPTZ
);
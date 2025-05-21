package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/models"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (u *Handlers) MockGenerateShortURL() string {
	return "short123"
}

type MockRepo struct {
	data map[string]string
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		data: make(map[string]string),
	}
}

// GetURL возвращает URL по токену
func (m *MockRepo) GetURL(ctx context.Context, token string) (string, error) {
	url, ok := m.data[token]
	if !ok {
		return "", sql.ErrNoRows
	}
	if url == "DELETED" {
		return "", sql.ErrNoRows
	}
	return url, nil
}

// GetURLWithDeleted — поддерживает флаг удаления
func (m *MockRepo) GetURLWithDeleted(ctx context.Context, token string) (string, bool, error) {
	url, ok := m.data[token]
	if !ok {
		return "", false, sql.ErrNoRows
	}
	if url == "DELETED" {
		return "", true, nil
	}
	return url, false, nil
}

// UpdateURL сохраняет URL (простейшая реализация)
func (m *MockRepo) UpdateURL(ctx context.Context, userID, shortURL, originalURL string) error {
	m.data[shortURL] = originalURL
	return nil
}

// BatchUpdateURL — сохраняет несколько URL для пользователя
func (m *MockRepo) BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error {
	for token, url := range URLs {
		m.data[token] = url
	}
	return nil
}

// GetAllURL — возвращает все URL пользователя
func (m *MockRepo) GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error) {
	var result []models.LinkPair
	for short, original := range m.data {
		if original == "DELETED" {
			continue
		}
		result = append(result, models.LinkPair{
			ShortURL: short,
			LongURL:  original,
		})
	}
	return result, nil
}

// Ping — проверка соединения (в моке всегда успешно)
func (m *MockRepo) Ping(ctx context.Context) error {
	return nil
}

func TestHandlers_PostHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name  string
		url   string
		wants want
	}{
		{
			name: "postTest",
			url:  "vk.com",
			wants: want{
				statusCode: 201,
			},
		},
		{
			name: "postTest2",
			url:  "vk.com/user/123",
			wants: want{
				statusCode: 201,
			},
		},
		{
			name: "postTest3",
			url:  "https://google.com/user/123",
			wants: want{
				statusCode: 201,
			},
		},
		{
			name: "postTest4",
			url:  "",
			wants: want{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			reqPost := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.url))
			ctx := utils.WithUserID(reqPost.Context(), "testUserID")
			reqPost = reqPost.WithContext(ctx)
			mockRepo := NewMockRepo()
			config := server.Config{}
			handler := NewHandlers(mockRepo, config)
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostHandler)
			h.ServeHTTP(rr, reqPost)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
		})
	}

}

func TestHandlers_GetHandler(t *testing.T) {
	type want struct {
		statusCode int
		location   string
		id         string
		response   string
	}

	tests := []struct {
		name      string
		id        string
		url       string
		setupMock func(*MockRepo)
		wants     want
	}{
		{
			name: "успешный редирект с http",
			id:   "abc123",
			url:  "http://vk.com",
			setupMock: func(repo *MockRepo) {
				repo.UpdateURL(utils.WithUserID(context.Background(), "user1"), "user1", "abc123", "http://vk.com")
			},
			wants: want{
				statusCode: 307,
				id:         "abc123",
				location:   "http://vk.com",
				response:   "",
			},
		},
		{
			name: "успешный редирект без схемы (добавится http://)",
			id:   "abc456",
			url:  "http://yandex.ru",
			setupMock: func(repo *MockRepo) {
				repo.UpdateURL(utils.WithUserID(context.Background(), "user1"), "user1", "abc456", "yandex.ru")
			},
			wants: want{
				statusCode: 307,
				id:         "abc456",
				location:   "http://yandex.ru",
				response:   "",
			},
		},
		{
			name:      "URL не найден — 404",
			id:        "abc345",
			url:       "",
			setupMock: func(repo *MockRepo) {}, // ничего не добавляем
			wants: want{
				statusCode: 404,
				id:         "",
				location:   "",
				response:   "404 page not found\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepo()
			config := server.Config{}
			handler := NewHandlers(mockRepo, config)
			reqGet := httptest.NewRequest(http.MethodGet, "http://localhost:8080/"+tt.wants.id, nil)
			ctx := utils.WithUserID(reqGet.Context(), "testUserID")
			reqGet = reqGet.WithContext(ctx)
			mockRepo.UpdateURL(ctx, "testUserID", tt.id, tt.url)
			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/{id}", handler.GetHandler)
			r.ServeHTTP(rr, reqGet)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
			assert.Equal(t, tt.wants.location, rr.Header().Get("Location"))
			if tt.wants.response != "" {
				assert.Equal(t, tt.wants.response, rr.Body.String())
			}
		})
	}
}

func TestHandlers_PostAPI(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name  string
		url   string
		wants want
	}{
		{
			name: "postTest",
			url:  `{"url": "vk.com"}`,
			wants: want{
				statusCode: 201,
			},
		},
		{
			name: "postTest2",
			url:  `{"url": "vk.com/user/123"}`,
			wants: want{
				statusCode: 201,
			},
		},
		{
			name: "postTest3",
			url:  `{"url": "https://google.com/user/123"}`,
			wants: want{
				statusCode: 201,
			},
		},
		{
			name: "postTest4",
			url:  "",
			wants: want{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			reqPost := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(tt.url))
			reqPost.Header.Set("Content-Type", "application/json")
			ctx := utils.WithUserID(reqPost.Context(), "testUserID")
			reqPost = reqPost.WithContext(ctx)
			repo := NewMockRepo()
			config := server.Config{}
			handle := NewHandlers(repo, config)
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(handle.PostAPI)
			h.ServeHTTP(rr, reqPost)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
		})
	}

}

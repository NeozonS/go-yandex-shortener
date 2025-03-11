package handlers

import (
	"bytes"
	"context"
	"errors"
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
	UpdateURLFunc      func(ctx context.Context, userID, shortURL, originURL string) error
	GetURLFunc         func(ctx context.Context, shortURL string) (string, error)
	GetAllURLFunc      func(ctx context.Context, userID string) ([]models.LinkPair, error)
	PingFunc           func(ctx context.Context) error
	BatchUpdateURLFunc func(ctx context.Context, userID string, URLs map[string]string) error
}

func (m *MockRepo) BatchUpdateURL(ctx context.Context, userID string, URLs map[string]string) error {
	if m.BatchUpdateURLFunc != nil {
		return m.BatchUpdateURLFunc(ctx, userID, URLs)
	}
	return nil
}

// UpdateURL реализует метод UpdateURL из интерфейса storage.Repository.
func (m *MockRepo) UpdateURL(ctx context.Context, userID, shortURL, originURL string) error {
	if m.UpdateURLFunc != nil {
		return m.UpdateURLFunc(ctx, userID, shortURL, originURL)
	}
	return nil
}

// GetURL реализует метод GetURL из интерфейса storage.Repository.
func (m *MockRepo) GetURL(ctx context.Context, shortURL string) (string, error) {
	if m.GetURLFunc != nil {
		return m.GetURLFunc(ctx, shortURL)
	}
	return "", nil
}

// GetAllURL реализует метод GetAllURL из интерфейса storage.Repository.
func (m *MockRepo) GetAllURL(ctx context.Context, userID string) ([]models.LinkPair, error) {
	if m.GetAllURLFunc != nil {
		return m.GetAllURLFunc(ctx, userID)
	}
	return []models.LinkPair{}, nil
}

func (m *MockRepo) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
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
			mockRepo := &MockRepo{}
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
		response   string
	}

	tests := []struct {
		name       string
		id         string
		mockGetURL func(ctx context.Context, shortURL string) (string, error)
		wants      want
	}{
		{
			name: "getTest - успешный запрос с http",
			id:   "abc123",
			mockGetURL: func(ctx context.Context, shortURL string) (string, error) {
				return "http://vk.com", nil
			},
			wants: want{
				statusCode: 307,
				location:   "http://vk.com",
				response:   "",
			},
		},
		{
			name: "getTest2 - успешный запрос без схемы",
			id:   "abc123",
			mockGetURL: func(ctx context.Context, shortURL string) (string, error) {
				return "vk.com", nil
			},
			wants: want{
				statusCode: 307,
				location:   "http://vk.com",
				response:   "",
			},
		},
		{
			name: "getTest3 - ошибка: URL не найден",
			id:   "abc123",
			mockGetURL: func(ctx context.Context, shortURL string) (string, error) {
				return "", errors.New("URL not found")
			},
			wants: want{
				statusCode: 400,
				location:   "",
				response:   "Запрашиваемая страница не найдена\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мок репозитория
			mockRepo := &MockRepo{
				GetURLFunc: tt.mockGetURL,
			}

			// Конфигурация
			config := server.Config{}

			// Создаем обработчик
			handler := NewHandlers(mockRepo, config)

			// Создаем запрос
			reqGet := httptest.NewRequest(http.MethodGet, "http://localhost:8080/"+tt.id, nil)

			// Добавляем userID в контекст (если нужно)
			ctx := utils.WithUserID(reqGet.Context(), "testUserID")
			reqGet = reqGet.WithContext(ctx)

			// Записываем ответ
			rr := httptest.NewRecorder()

			// Используем chi для маршрутизации
			r := chi.NewRouter()
			r.Get("/{id}", handler.GetHandler)
			r.ServeHTTP(rr, reqGet)

			// Проверяем статус код
			assert.Equal(t, tt.wants.statusCode, rr.Code)

			// Проверяем заголовок Location
			assert.Equal(t, tt.wants.location, rr.Header().Get("Location"))

			// Проверяем тело ответа (для ошибок)
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
			repo := &MockRepo{}
			config := server.Config{}
			handle := NewHandlers(repo, config)
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(handle.PostAPI)
			h.ServeHTTP(rr, reqPost)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
		})
	}

}

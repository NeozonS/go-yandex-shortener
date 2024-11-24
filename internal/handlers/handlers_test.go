package handlers

import (
	"bytes"
	"github.com/NeozonS/go-shortener-ya.git/internal/server"
	"github.com/NeozonS/go-shortener-ya.git/internal/storage/mapbd"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (u *Handlers) MockgenerateShortURL() string {
	return "short123"
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
			config := server.Config{}
			repo := mapbd.New()
			handl := NewHandlers(repo, config)
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(handl.PostHandler)
			h.ServeHTTP(rr, reqPost)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
		})
	}

}

func TestHandlers_GetHandler(t *testing.T) {
	type want struct {
		statusCode int
		url        string
	}

	tests := []struct {
		name  string
		url   string
		badid string
		wants want
	}{
		{
			name: "getTest",
			url:  "http://vk.com",
			wants: want{
				statusCode: 307,
				url:        "http://vk.com",
			},
		},
		{
			name: "getTest2",
			url:  "vk.com",
			wants: want{
				statusCode: 307,
				url:        "http://vk.com",
			},
		},
		{
			name: "getTest3",
			url:  "google.com",
			wants: want{
				statusCode: 307,
				url:        "http://google.com",
			},
		},
		{
			name: "getTest4",
			url:  "",
			wants: want{
				statusCode: 400,
				url:        "",
			},
		},
		{
			name:  "getTest5",
			url:   "google.com",
			badid: "123",
			wants: want{
				statusCode: 400,
				url:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mapbd.New()
			config := server.Config{}
			handl := NewHandlers(repo, config)
			id := handl.generateShortURL()
			repo.UpdateURL(tt.url, id)
			reqGet := httptest.NewRequest(http.MethodGet, "http://localhost:8080/"+id+tt.badid, nil)
			r := chi.NewRouter()
			r.Get("/{id}", handl.GetHandler)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, reqGet)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
			assert.Equal(t, tt.wants.url, rr.Header().Get("Location"))

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
			repo := mapbd.New()
			config := server.Config{}
			handl := NewHandlers(repo, config)
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(handl.PostAPI)
			h.ServeHTTP(rr, reqPost)
			assert.Equal(t, tt.wants.statusCode, rr.Code)
		})
	}

}

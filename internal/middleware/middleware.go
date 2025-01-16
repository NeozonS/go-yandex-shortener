package middleware

import (
	"compress/gzip"
	"context"
	"github.com/NeozonS/go-shortener-ya.git/internal/utils"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (gz gzipResponseWriter) Write(data []byte) (int, error) {
	return gz.Writer.Write(data)
}

func GzipRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzRead, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Unable to decode gzip body", http.StatusBadRequest)
				return
			}
			defer gzRead.Close()
			r.Body = gzRead
		}
		next.ServeHTTP(w, r)
	})
}

func GzipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		gz := gzipResponseWriter{Writer: gzWriter, ResponseWriter: w}
		next.ServeHTTP(gz, r)
	})
}

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromCookie(r)
		if err != nil {
			log.Printf("Error getting userID from cookie: %v", err)
			http.Error(w, "Failed to decrypt userID", http.StatusBadRequest)
			return
		}

		if userID == "" {
			newUserID := uuid.New().String()
			utils.SetCookie(w, newUserID)
			ctx := context.WithValue(r.Context(), userIDKey, newUserID)
			log.Printf("New userID generated: %s", newUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		log.Printf("UserID found in cookie: %s", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

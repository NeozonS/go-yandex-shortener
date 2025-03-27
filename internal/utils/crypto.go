package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var secretKey = []byte("temporary-secret-key-for-development")

func Encrypt(data string) (string, error) {
	key := sha256.Sum256(secretKey)

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// Генерируем уникальный nonce длиной 12 байт
	nonce := key[len(key)-aesGCM.NonceSize():]

	ciphertext := aesGCM.Seal(nil, nonce, []byte(data), nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(data string) (string, error) {
	key := sha256.Sum256(secretKey)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := key[len(key)-aesGCM.NonceSize():]
	decodedData, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	plaintext, err := aesGCM.Open(nil, nonce, decodedData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func GenerateShortURL(originalURL, userID string) (string, error) {
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes, %w", err)
	}
	data := fmt.Sprintf("%s%s%x", originalURL, userID, randomBytes)
	hash := sha256.Sum256([]byte(data))
	token := base64.URLEncoding.EncodeToString(hash[:])
	return token[:8], nil
}

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func GenerateShortURL(originalURL, userID string) string {
	hash := sha256.Sum256([]byte(originalURL + userID))
	token := base64.URLEncoding.EncodeToString(hash[:])
	return token[:8]
}

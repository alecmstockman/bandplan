package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateSessionExpiration() time.Time {
	expiration := time.Now().Add(1 * time.Hour)
	return expiration
}

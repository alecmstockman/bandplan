package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func GenerateSessionToken() (string, error) {
	log.Println("- Helper - GenerateSessionToken")
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"bandplan/src/database"
	"bandplan/src/models"
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

func GetAuthenticatedUser(r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return models.User{}, err
	}
	fmt.Println("Get Authenticated User cookie: ", cookie)

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	if err != nil {
		fmt.Println("Unable to get user by token: ", err)
		return models.User{}, nil
	}

	return user, nil
}

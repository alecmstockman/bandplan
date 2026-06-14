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

func HelperGenerateSessionToken() (string, error) {
	fmt.Println(" - HelperGenerateSessionToken")
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HelperGenerateSessionExpiration() time.Time {
	fmt.Println(" - HelperGenerateSessionExpiration")
	expiration := time.Now().Add(1 * time.Hour)
	return expiration
}

func HelperGetAuthenticatedUser(r *http.Request) (models.User, error) {
	fmt.Println(" - HelperGetAuthenticatedUser")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return models.User{}, err
	}
	fmt.Println("cookie: ", cookie)

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	fmt.Println("user: ", user)
	if err != nil {
		fmt.Println("Unable to get user by token: ", err)
		return models.User{}, err
	}

	return user, nil
}

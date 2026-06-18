package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"bandplan/src/database"
	"bandplan/src/models"
)

func HelperGenerateSessionToken() (string, error) {
	log.Println("- HelperGenerateSessionToken")
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ProcessBandNameEntry(bandNameEntry string) string {
	log.Println("- ProcessBandNameEntry called")
	stripped := strings.TrimSpace(bandNameEntry)
	cleanName := strings.ToLower(stripped)
	return cleanName
}

func HelperGenerateSessionExpiration() time.Time {
	log.Println("- HelperGenerateSessionExpiration")
	expiration := time.Now().Add(1 * time.Hour)
	return expiration
}

func HelperGetAuthenticatedUser(r *http.Request) (models.User, error) {
	log.Println("- HelperGetAuthenticatedUser")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return models.User{}, err
	}

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	fmt.Printf("-------- helper get auth user: %+v\n", user)
	if err != nil {
		fmt.Println(" - Unable to get user by token: ", err)
		return models.User{}, err
	}

	return user, nil
}

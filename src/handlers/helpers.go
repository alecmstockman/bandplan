package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"bandplan/src/database"
	"bandplan/src/models"

	"github.com/google/uuid"
)

func HelperGetAuthContext(r *http.Request) (AuthContext, error) {
	log.Println("- HelperGetAuthContext")

	auth, ok := r.Context().Value(AuthContextKey).(AuthContext)
	if !ok {
		log.Println("   auth OK status: ", ok)
		return AuthContext{}, errors.New("auth context missing from request")
	}
	return auth, nil
}

func HelperGetAuthenticatedUserAndBand(r *http.Request) (models.User, models.Band, error) {
	fmt.Println("----------------------------------------------")
	log.Println("- HelperGetAuthenticatedUserAndBand")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return models.User{}, models.Band{}, err
	}

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	if err != nil {
		log.Println("   Unable to get users by token: ", err)
		return models.User{}, models.Band{}, err
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		log.Println("   Unable to get band by user ID: ", err)
		return models.User{}, models.Band{}, err
	}

	// session, err := database.SessionsTableGetSessionByToken(cookie.Value)
	// if err != nil {
	// 	log.Println("   Unable to get session by token: ", err)
	// 	return models.User{}, models.Band{}, err
	// }
	// fmt.Println("Session: ", session)

	// user, band, err := database.SessionsTableGetAuthContextByToken(cookie.Value)
	// if err != nil {
	// 	log.Println("   Unable to get user and band by token: ", err)
	// 	return models.User{}, models.Band{}, err
	// }

	return user, band, nil
}

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

func HelperMakeUserSlug(userName string) string {
	name := strings.ToLower(strings.TrimSpace(userName))

	var b strings.Builder
	lastDash := false

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	id := "-" + uuid.NewString()[:6]

	return strings.Trim(b.String(), "-") + id
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
		log.Println("   Unable to get session token from cookie", err)
		return models.User{}, err
	}

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	if err != nil {
		log.Println("   Error getting user from sessions table by token: ", err)
		return models.User{}, err
	}

	return user, nil
}

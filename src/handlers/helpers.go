package handlers

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
)

func HelperGetAuthContext(r *http.Request) (AuthContext, error) {

	auth, ok := r.Context().Value(AuthContextKey).(AuthContext)
	if !ok {
		return AuthContext{}, errors.New("auth context missing from request")
	}
	return auth, nil
}

func HelperGetAuthenticatedUserAndBand(r *http.Request) (models.User, models.Band, error) {

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

	return user, band, nil
}

func HelperProcessBandNameEntry(bandNameEntry string) string {

	stripped := strings.TrimSpace(bandNameEntry)
	cleanName := strings.ToLower(stripped)
	return cleanName
}

func HelperGenerateSessionExpiration() time.Time {
	expiration := time.Now().Add(1 * time.Hour)
	return expiration
}

func HelperGetAuthenticatedUser(r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		slog.Error(
			"unable to get session token",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		return models.User{}, err
	}

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	if err != nil {
		slog.Error(
			"unable to get session by token",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		return models.User{}, err
	}

	return user, nil
}

func HelperITunesArtworkURLLarge(url string) string {
	return ""
}

func FormatOptionalTime(value *time.Time, timezone string) string {
	if value == nil {
		return ""
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = time.UTC
	}

	return value.In(location).Format("3:04 PM")
}

func ValidatePriceEntry(price string) bool {
	fmt.Println("price: ", price)

	decimalCount := 0

	for _, r := range price {
		if r == '.' && len(price) == 1 {
			return false
		}
		if r == '.' {
			decimalCount++

			if decimalCount > 1 {
				return false
			}
			continue
		}

		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

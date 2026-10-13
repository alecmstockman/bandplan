package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"bandplan/src/database"
	"bandplan/src/helpers"
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

	tokenHash := helpers.HashSessionToken(cookie.Value)

	user, err := database.SessionsTableGetUserByToken(tokenHash)
	if err != nil {
		log.Println("   Unable to get users by token: ", err)
		return models.User{}, models.Band{}, err
	}

	band, err := database.SessionsTableGetCurrentBandByToken(tokenHash, user.UserID)
	if err == nil {
		return user, band, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Println("   Unable to get current band from session: ", err)
		return models.User{}, models.Band{}, err
	}

	band, err = database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		log.Println("   Unable to get fallback band by user ID: ", err)
		return models.User{}, models.Band{}, err
	}

	updated, err := database.SessionsTableSetCurrentBand(tokenHash, user.UserID, band.BandID)
	if err != nil {
		return models.User{}, models.Band{}, err
	}
	if !updated {
		return models.User{}, models.Band{}, fmt.Errorf("unable to set current band for active session")
	}

	return user, band, nil
}

func HelperGenerateSessionExpiration() time.Time {
	expiration := time.Now().Add(1 * time.Hour)
	return expiration
}

func HelperGetAuthenticatedUser(r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("   ERROR: HelperGetAuthenticatedUser, SessionsTableGetUserByToken: ", err)
		// slog.Error(
		// 	"unable to get session token",
		// 	"request_id", requestlog.GetRequestID(r.Context()),
		// 	"error", err,
		// )
		return models.User{}, err
	}

	tokenHash := helpers.HashSessionToken(cookie.Value)

	user, err := database.SessionsTableGetUserByToken(tokenHash)
	if err != nil {
		slog.Error(
			"unable to get session by token",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
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

func parseOptionalDueFields(dateValue, timeValue, timezoneValue string) (*time.Time, *time.Time, *string, error) {
	if dateValue == "" && timeValue == "" {
		return nil, nil, nil, nil
	}

	location, err := time.LoadLocation(timezoneValue)
	if err != nil {
		return nil, nil, nil, errors.New("Invalid due timezone")
	}

	var dueDate *time.Time
	if dateValue != "" {
		parsedDate, err := time.ParseInLocation("2006-01-02", dateValue, location)
		if err != nil {
			return nil, nil, nil, errors.New("Invalid due date")
		}
		dueDate = &parsedDate
	}

	var dueTime *time.Time
	if timeValue != "" {
		parsedTime, err := time.ParseInLocation("15:04", timeValue, location)
		if err != nil {
			return nil, nil, nil, errors.New("Invalid due time")
		}
		dueTime = &parsedTime
	}

	return dueDate, dueTime, &timezoneValue, nil
}

func ValidatePriceEntry(price string) bool {

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

func UpcomingEvents(events []models.Event) []models.Event {
	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	var upcomingEvents []models.Event

	for _, event := range events {
		eventDate := time.Date(
			event.EventDate.Year(),
			event.EventDate.Month(),
			event.EventDate.Day(),
			0,
			0,
			0,
			0,
			today.Location(),
		)

		if !eventDate.Before(today) {
			upcomingEvents = append(upcomingEvents, event)
		}
	}
	return upcomingEvents
}

func PastEvents(events []models.Event) []models.Event {
	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	var pastEvents []models.Event

	for _, event := range events {
		eventDate := time.Date(
			event.EventDate.Year(),
			event.EventDate.Month(),
			event.EventDate.Day(),
			0,
			0,
			0,
			0,
			today.Location(),
		)

		if eventDate.Before(today) {
			pastEvents = append(pastEvents, event)
		}
	}
	return pastEvents
}

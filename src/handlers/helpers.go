package handlers

import (
	"errors"
	"log"
	"log/slog"
	"net/http"

	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"bandplan/src/services"
)

func writeImageUploadError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrUnsupportedImageType):
		http.Error(w, "Unsupported image type", http.StatusUnsupportedMediaType)
	case errors.Is(err, services.ErrImageDimensionsExceeded):
		http.Error(w, "Image dimensions exceed upload limits", http.StatusRequestEntityTooLarge)
	case errors.Is(err, services.ErrInvalidImage):
		http.Error(w, "Invalid image", http.StatusBadRequest)
	default:
		http.Error(w, "Unable to process image", http.StatusInternalServerError)
	}
}

func HelperGetAuthContext(r *http.Request) (AuthContext, error) {

	auth, ok := r.Context().Value(AuthContextKey).(AuthContext)
	if !ok {
		return AuthContext{}, errors.New("auth context missing from request")
	}
	return auth, nil
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

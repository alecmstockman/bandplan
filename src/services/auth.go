package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alexedwards/argon2id"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s Service) LoginValidation(ctx context.Context, email, password string) (models.Session, error) {
	if email == "" || len(email) > 254 {
		return models.Session{}, ErrInvalidCredentials
	}

	normalizedEmail := helpers.NormalizeEmail(email)
	validatedEmail := helpers.ValidateEmail(normalizedEmail)

	if !validatedEmail {
		return models.Session{}, ErrInvalidCredentials
	}

	if len(password) < 8 || len(password) > 255 {
		return models.Session{}, ErrInvalidCredentials
	}

	user, err := database.UsersTableGetUserByEmail(normalizedEmail)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := helpers.RunDummyArgon2Hash(password)
		if err != nil {
			return models.Session{}, fmt.Errorf("compare dummy password hash: %w", err)
		}
		return models.Session{}, ErrInvalidCredentials
	}
	if err != nil {
		return models.Session{}, fmt.Errorf("get user by email: %w", err)
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return models.Session{}, fmt.Errorf("compare password hash: %w", err)

	}
	if !match {
		return models.Session{}, ErrInvalidCredentials
	}

	token, err := helpers.GenerateSessionToken()
	if err != nil {
		slog.Error(
			"unable to generate session token",
			"request_id", requestlog.GetRequestID(ctx),
			"error", err,
		)
		return models.Session{}, fmt.Errorf("unable to generate session token: %w", err)
	}

	params := models.CreateSessionParams{
		UserID: user.UserID,
		Token:  token,
	}

	session, err := database.SessionsTableCreateSession(params)
	if err != nil {
		return models.Session{}, fmt.Errorf("unable to create session: %w", err)
	}

	return session, nil

}

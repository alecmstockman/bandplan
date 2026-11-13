package services

import (
	"bandplan/src/clients"
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alexedwards/argon2id"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLoginRateLimited   = errors.New("login rate limited")
	ErrNotVerified        = errors.New("user not verified")
)

type LoginLimiter interface {
	Allow(string) bool
	Reset(string)
}

func (s Service) LoginValidation(ctx context.Context, email, password string) (models.Session, error) {
	if email == "" || len(email) > 254 {
		return models.Session{}, ErrInvalidCredentials
	}

	normalizedEmail := helpers.NormalizeEmail(email)
	validatedEmail := helpers.ValidateEmail(normalizedEmail)

	if !validatedEmail {
		return models.Session{}, ErrInvalidCredentials
	}

	if !s.LoginEmailLimiter.Allow(normalizedEmail) {
		return models.Session{}, ErrLoginRateLimited
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

	verified, err := database.UsersTableCheckEmailVerification(user.UserID)
	if err != nil {
		return models.Session{}, err
	}
	if verified != true {
		return models.Session{}, ErrNotVerified
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

	tokenHash := helpers.HashSessionToken(token)

	params := models.CreateSessionParams{
		UserID:    user.UserID,
		Token:     token,
		TokenHash: tokenHash,
	}

	session, err := database.SessionsTableCreateSession(params)
	if err != nil {
		return models.Session{}, fmt.Errorf("unable to create session: %w", err)
	}

	return session, nil

}

func (s Service) PasswordResetSendEmail(ctx context.Context, email string) error {

	if email == "" || len(email) > 254 {
		return ErrInvalidCredentials
	}

	normalizedEmail := helpers.NormalizeEmail(email)
	validatedEmail := helpers.ValidateEmail(normalizedEmail)

	if !validatedEmail {
		return ErrInvalidCredentials
	}

	if !s.LoginEmailLimiter.Allow(normalizedEmail) {
		return ErrLoginRateLimited
	}

	user, err := database.UsersTableGetUserByEmail(normalizedEmail)
	if err != nil {
		return fmt.Errorf("Unable to get user")
	}

	verified, err := database.UsersTableCheckEmailVerification(user.UserID)
	if err != nil {
		return err
	}
	if verified != true {
		return ErrNotVerified
	}

	token, err := helpers.GenerateSessionToken()
	if err != nil {
		return err
	}

	tokenHash := helpers.HashRegistrationCode(token)

	expiration := time.Now().Add(30 * time.Minute)

	err = database.PasswordResetTokensTableCreateToken(ctx, user.UserID, tokenHash, expiration)
	if err != nil {
		return err
	}

	data := models.PasswordResetEmailData{
		Token: token,
		Email: user.Email,
	}

	err = clients.SendPasswordResetEmail(data)
	if err != nil {
		return err
	}

	return nil

}

func (s Service) PasswordReset(ctx context.Context, token, password string) error {

	if len(token) != 64 {
		return errors.New("invalid reset token")
	}

	if _, err := hex.DecodeString(token); err != nil {
		return errors.New("invalid reset token")
	}

	tokenHash := helpers.HashRegistrationCode(token)

	if len(password) < 8 || len(password) > 255 {
		return ErrInvalidCredentials
	}

	valid := helpers.PasswordValidateLength(password)
	if valid != true {
		return errors.New("invalid password length")
	}

	passwordHash, err := helpers.HashPassword(password)
	if err != nil {
		return errors.New("unable to hash password")
	}

	_, err = database.PasswordResetTokenTableValidateCodeReturnUserID(ctx, tokenHash, passwordHash)
	if err != nil {
		return err
	}

	return nil

}

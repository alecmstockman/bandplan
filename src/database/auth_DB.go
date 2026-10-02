package database

import (
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func RegisterInitialUserBandAndChat(ctx context.Context, user models.User, band models.Band, registrationTokenHash string) error {
	log.Println("RegisterInitialUserBandAndChat")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	userQuery := `
		INSERT INTO users (
			user_id,
			name,
			first_name,
			last_name,
			display_name,
			email,
			slug,
			password_hash,
			is_admin,
			legal_accepted,
			legal_accepted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		`

	_, err = tx.Exec(
		userQuery,
		user.UserID,
		user.Name,
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.Email,
		user.Slug,
		user.PasswordHash,
		user.IsAdmin,
		user.LegalAccepted,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	bandQuery := `
		INSERT INTO bands(
			band_id,
			name,
			slug,
			created_by,
			updated_by
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(
		bandQuery,
		band.BandID,
		band.Name,
		band.Slug,
		user.UserID,
		user.UserID,
	)

	if err != nil {
		return fmt.Errorf("create band: %w", err)
	}

	membersQuery := `
		INSERT INTO band_members(
			band_id,
			user_id
		) VALUES ($1, $2)
	`

	_, err = tx.Exec(
		membersQuery,
		band.BandID,
		user.UserID,
	)

	if err != nil {
		return fmt.Errorf("insert band member: %w", err)
	}

	chatID := uuid.New().String()
	chatName := fmt.Sprintf("%s (Band Chat)", band.Name)

	primaryChatQuery := `
		INSERT INTO chats (
			chat_id,
			band_id,
			name,
			slug,
			is_primary,
			created_by,
			updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err = tx.Exec(
		primaryChatQuery,
		chatID,
		band.BandID,
		chatName,
		band.Slug,
		true,
		user.UserID,
		user.UserID,
	)

	if err != nil {
		return fmt.Errorf("create primary chat: %w", err)
	}

	chatMembersQuery := `
		INSERT INTO chat_members(
			chat_id,
			user_id
		) VALUES (
			$1, $2
		)
	`

	_, err = tx.Exec(
		chatMembersQuery,
		chatID,
		user.UserID,
	)
	if err != nil {
		return fmt.Errorf("insert chat member: %w", err)
	}

	registrationCodeQuery := `
		DELETE FROM user_registrations
		WHERE user_registration_id = $1
			AND email = $2
			AND expires_at > NOW()
	`

	result, err := tx.ExecContext(
		ctx,
		registrationCodeQuery,
		registrationTokenHash,
		user.Email,
	)
	if err != nil {
		return fmt.Errorf("delete registration code: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count deleted registration codes: %w", err)
	}

	if count == 0 {
		return errors.New("invalid or expired registration code")
	}
	if count > 1 {
		return errors.New("registration code matches multiple active rows")
	}

	return tx.Commit()
}

func RegisterNewBandUser(ctx context.Context, user models.User, bandID, chatID, registrationTokenHash, accessCodeHash string) (string, error) {
	log.Println("- RegisterNewBandUser")

	tx, err := DB.Begin()
	if err != nil {
		return "", err
	}

	defer tx.Rollback()

	userQuery := `
		INSERT INTO users (
			user_id,
			name,
			first_name,
			last_name,
			display_name,
			email,
			slug,
			password_hash,
			is_admin,
			legal_accepted,
			legal_accepted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		`

	_, err = tx.Exec(
		userQuery,
		user.UserID,
		user.Name,
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.Email,
		user.Slug,
		user.PasswordHash,
		user.IsAdmin,
		user.LegalAccepted,
		time.Now(),
	)

	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	registrationCodeQuery := `
		DELETE FROM user_registrations
		WHERE user_registration_id = $1
			AND email = $2
			AND expires_at > NOW()
	`

	result, err := tx.ExecContext(
		ctx,
		registrationCodeQuery,
		registrationTokenHash,
		user.Email,
	)
	if err != nil {
		return "", fmt.Errorf("delete registration code: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("count deleted registration codes: %w", err)
	}

	if count == 0 {
		return "", errors.New("invalid or expired registration code")
	}
	if count > 1 {
		return "", errors.New("registration code matches multiple active rows")
	}

	accessCodeQuery := `
		DELETE FROM access_codes
		WHERE code_hash = $1
		AND band_id = $2
		AND expires_at > NOW()
	`

	result, err = tx.ExecContext(
		ctx,
		accessCodeQuery,
		accessCodeHash,
		bandID,
	)

	if err != nil {
		return "", fmt.Errorf("delete access code: %w", err)
	}

	count, err = result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("count deleted access codes: %w", err)
	}

	if count == 0 {
		return "", errors.New("invalid or expired access code")
	}
	if count > 1 {
		return "", errors.New("access code matches multiple active rows")
	}

	membersQuery := `
		INSERT INTO band_members(
			band_id,
			user_id
		) VALUES ($1, $2)
	`

	_, err = tx.Exec(
		membersQuery,
		bandID,
		user.UserID,
	)

	if err != nil {
		return "", fmt.Errorf("insert band member: %w", err)
	}

	chatMembersQuery := `
		INSERT INTO chat_members(
			chat_id,
			user_id
		) VALUES (
			$1, $2
		)
	`

	_, err = tx.Exec(
		chatMembersQuery,
		chatID,
		user.UserID,
	)
	if err != nil {
		return "", fmt.Errorf("insert chat member: %w", err)
	}

	return bandID, tx.Commit()
}

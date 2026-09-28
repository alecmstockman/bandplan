package database

import (
	"bandplan/src/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func RegisterInitialUserBandAndChat(user models.User, band models.Band, registrationID string) error {
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
			is_admin
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
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

	registrationIDQuery := `
		DELETE FROM user_registrations
		WHERE user_registration_id = $1
	`
	_, err = tx.Exec(
		registrationIDQuery,
		registrationID,
	)
	if err != nil {
		return fmt.Errorf("delete user registration: %w", err)
	}

	return tx.Commit()
}

func RegisterNewBandUser(user models.User, bandID, chatID, accessCode string) error {
	log.Println("- RegisterNewBandUser")

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
			is_admin
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
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
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
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
		return fmt.Errorf("insert band member: %w", err)
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

	hash := sha256.Sum256([]byte(accessCode))
	codeHash := hex.EncodeToString(hash[:])

	deleteAccessCodeQuery := `
		DELETE FROM access_codes
		WHERE code_hash = $1
			AND band_id = $2

	`

	result, err := tx.Exec(
		deleteAccessCodeQuery,
		codeHash,
		bandID,
	)
	if err != nil {
		return fmt.Errorf("delete access code: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.New("No rows deleted")
	}
	if rowsAffected > 1 {
		return errors.New("Unable to delete, multiple matches found")
	}

	return tx.Commit()
}

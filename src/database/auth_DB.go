package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func RegisterNewUserBandAndChat(user models.User, band models.Band) error {
	fmt.Println("----------------------------------------------")
	log.Println("RegisterNewUserBandAndChat")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	userQuery := `
		INSERT INTO users (
			user_id,
			name,
			display_name,
			email,
			slug,
			password_hash,
			is_admin
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		`

	_, err = tx.Exec(
		userQuery,
		user.UserID,
		user.Name,
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

	return tx.Commit()
}

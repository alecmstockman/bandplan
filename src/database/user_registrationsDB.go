package database

import (
	"bandplan/src/models"
	"fmt"
	"log"
)

func UsersRegTableCreateInititialUser(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("UsersRegTableCreateInititialUser")

	query := `
		INSERT INTO user_registrations (
			user_registration_id,
			access_code_hash,
			first_name,
			last_name,
			display_name,
			timezone,
			email
		) VALUES (
			$1, NULLIF($2, ''), $3, $4, $5, $6, $7
		) 
		RETURNING
			id,
			user_registration_id,
			COALESCE(access_code_hash, ''),
			first_name,
			last_name,
			display_name,
			timezone,
			email,
			email_verified,
			created_at,
			updated_at,
			expires_at
	`

	var newUser models.UserRegistration

	err := DB.QueryRow(
		query,
		user.UserRegistrationID,
		user.AccessCodeHash,
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.Timezone,
		user.Email,
	).Scan(
		&newUser.ID,
		&newUser.UserRegistrationID,
		&newUser.AccessCodeHash,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.DisplayName,
		&newUser.Timezone,
		&newUser.Email,
		&newUser.EmailVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
		&newUser.ExpiresAt,
	)
	if err != nil {
		return models.UserRegistration{}, err
	}

	return newUser, nil
}

func UsersRegTableUpdateInitialUser(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("UsersRegTableUpdateInitialUser")

	query := `
		UPDATE user_registrations
		SET
			access_code_hash = NULLIF($2, ''),
			first_name = $3,
			last_name = $4,
			display_name = $5,
			timezone = $6,
			email = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_registration_id = $1
		RETURNING
			id,
			user_registration_id,
			COALESCE(access_code_hash, ''),
			first_name,
			last_name,
			display_name,
			timezone,
			email,
			email_verified,
			created_at,
			updated_at,
			expires_at
	`

	var updatedUser models.UserRegistration

	err := DB.QueryRow(
		query,
		user.UserRegistrationID,
		user.AccessCodeHash,
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.Timezone,
		user.Email,
	).Scan(
		&updatedUser.ID,
		&updatedUser.UserRegistrationID,
		&updatedUser.AccessCodeHash,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.DisplayName,
		&updatedUser.Timezone,
		&updatedUser.Email,
		&updatedUser.EmailVerified,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
		&updatedUser.ExpiresAt,
	)
	if err != nil {
		return models.UserRegistration{}, err
	}

	return updatedUser, nil
}

func UsersRegTableGetUserByID(userID string) (models.UserRegistration, error) {
	log.Println("- UsersRegTableGetUserByID")

	query := `
		SELECT 
			id,
			user_registration_id,
			COALESCE(access_code_hash, ''),
			first_name,
			last_name,
			display_name,
			timezone,
			email,
			email_verified,
			created_at,
			updated_at,
			expires_at
		FROM user_registrations
		WHERE user_registration_id = $1
	`

	var newUser models.UserRegistration

	err := DB.QueryRow(query, userID).Scan(
		&newUser.ID,
		&newUser.UserRegistrationID,
		&newUser.AccessCodeHash,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.DisplayName,
		&newUser.Timezone,
		&newUser.Email,
		&newUser.EmailVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
		&newUser.ExpiresAt,
	)
	if err != nil {
		return models.UserRegistration{}, err
	}

	return newUser, nil
}

func UsersRegTableDeleteUserByID(registrationID string) error {
	log.Println("- UsersRegTableDeleteUserByID")

	query := `
		DELETE FROM user_registrations
		WHERE user_registration_id = $1
	`

	err := DB.QueryRow(query, registrationID)
	if err != nil {
		log.Println("Unable to delete user registration", err)
		return fmt.Errorf("err: %v", err)
	}

	return nil
}

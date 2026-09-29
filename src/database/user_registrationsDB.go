package database

import (
	"bandplan/src/models"
	"errors"
	"fmt"
	"log"
	"time"
)

func UsersRegTableCreateInititialUser(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("UsersRegTableCreateInititialUser")

	expiresAt := time.Now().Add(10 * time.Minute)

	query := `
		INSERT INTO user_registrations (
			user_registration_id,
			access_code_hash,
			first_name,
			last_name,
			display_name,
			band_id,
			timezone,
			email,
			expires_at
		) VALUES (
			$1, NULLIF($2, ''), $3, $4, $5, NULLIF($6, ''), $7, $8, $9
		) 
		RETURNING
			id,
			user_registration_id,
			COALESCE(access_code_hash, ''),
			first_name,
			last_name,
			display_name,
			COALESCE(band_id, ''),
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
		user.BandID,
		user.Timezone,
		user.Email,
		expiresAt,
	).Scan(
		&newUser.ID,
		&newUser.UserRegistrationID,
		&newUser.AccessCodeHash,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.DisplayName,
		&newUser.BandID,
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
			band_id = NULLIF($6, ''), 
			timezone = $7,
			email = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_registration_id = $1
			AND expires_at > NOW()
		RETURNING
			id,
			user_registration_id,
			COALESCE(access_code_hash, ''),
			first_name,
			last_name,
			display_name,
			COALESCE(band_id, ''),
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
		user.BandID,
		user.Timezone,
		user.Email,
	).Scan(
		&updatedUser.ID,
		&updatedUser.UserRegistrationID,
		&updatedUser.AccessCodeHash,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.DisplayName,
		&updatedUser.BandID,
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

	if updatedUser.ExpiresAt.Before(time.Now().UTC()) {
		return models.UserRegistration{}, errors.New("inalid registration id")
	}

	return updatedUser, nil
}

func UserRegTableValidateRegistrationID(registrationID string) (bool, error) {
	log.Println("- UserRegTableValidateRegistrationID")

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_registrations
			WHERE user_registration_id = $1
				AND expires_at > NOW()
		)
	`
	var valid bool
	err := DB.QueryRow(
		query,
		registrationID,
	).Scan(
		&valid,
	)

	if err != nil {
		return false, err
	}
	return valid, nil
}

func UsersRegTableGetUserRegistrationID(userID string) (models.UserRegistration, error) {
	log.Println("- UsersRegTableGetUserRegistrationID")

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
			AND expires_at > NOW()
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

	if newUser.ExpiresAt.Before(time.Now().UTC()) {
		return models.UserRegistration{}, errors.New("inalid registration id")
	}

	return newUser, nil
}

func UsersRegTableDeleteUserByID(registrationID string) (bool, error) {
	log.Println("- UsersRegTableDeleteUserByID")

	query := `
		DELETE FROM user_registrations
		WHERE user_registration_id = $1
	`

	result, err := DB.Exec(query, registrationID)
	if err != nil {
		log.Println("Unable to delete user registration", err)
		return false, fmt.Errorf("err: %v", err)
	}

	deleted, err := result.RowsAffected()

	return deleted > 0, nil
}

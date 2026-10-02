package database

import (
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var ErrRegistrationExpired = errors.New("registration expired")

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
		user.RegistrationTokenHash,
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
		&newUser.RegistrationTokenHash,
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
		user.RegistrationTokenHash,
		user.AccessCodeHash,
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.BandID,
		user.Timezone,
		user.Email,
	).Scan(
		&updatedUser.ID,
		&updatedUser.RegistrationTokenHash,
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserRegistration{}, ErrRegistrationExpired
		}
		return models.UserRegistration{}, err
	}

	if updatedUser.ExpiresAt.Before(time.Now().UTC()) {
		return models.UserRegistration{}, ErrRegistrationExpired
	}

	return updatedUser, nil
}

func UserRegTableValidateRegistrationToken(registrationTokenHash string) (bool, error) {
	log.Println("- UserRegTableValidateRegistrationToken")

	query := `
		SELECT expires_at
		FROM user_registrations
		WHERE user_registration_id = $1
	`
	var expiresAt time.Time
	err := DB.QueryRow(
		query,
		registrationTokenHash,
	).Scan(
		&expiresAt,
	)

	if err != nil {
		return false, err
	}
	if !expiresAt.After(time.Now().UTC()) {
		return false, ErrRegistrationExpired
	}
	return true, nil
}

func UsersRegTableGetRegistrationToken(registrationTokenHash string) (models.UserRegistration, error) {
	log.Println("- UsersRegTableGetRegistrationToken")

	query := `
		SELECT 
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
		FROM user_registrations
		WHERE user_registration_id = $1
	`

	var newUser models.UserRegistration

	err := DB.QueryRow(query, registrationTokenHash).Scan(
		&newUser.ID,
		&newUser.RegistrationTokenHash,
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

	if newUser.ExpiresAt.Before(time.Now().UTC()) {
		return models.UserRegistration{}, ErrRegistrationExpired
	}

	return newUser, nil
}

func UsersRegTableDeleteUserByID(registrationTokenHash string) (bool, error) {
	log.Println("- UsersRegTableDeleteUserByID")

	query := `
		DELETE FROM user_registrations
		WHERE user_registration_id = $1
	`

	result, err := DB.Exec(query, registrationTokenHash)
	if err != nil {
		log.Println("Unable to delete user registration", err)
		return false, fmt.Errorf("err: %v", err)
	}

	deleted, err := result.RowsAffected()

	return deleted > 0, nil
}

func UsersRegTableDeleteExpiredUserByEmail(email string) error {
	log.Println("- UsersRegTableDeleteExpiredUserByEmail")

	query := `
		DELETE FROM user_registrations
		WHERE LOWER(email) = LOWER($1)
			AND expires_at <= NOW()
	`

	_, err := DB.Exec(query, email)
	if err != nil {
		log.Println("Unable to delete user registration", err)
		return fmt.Errorf("err: %v", err)
	}

	return nil
}

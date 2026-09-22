package database

import (
	"bandplan/src/models"
	"log"
)

func UsersRegTableCreateInititialUser(user models.UserRegistration) (models.UserRegistration, error) {
	log.Println("UsersRegTableCreateInititialUser")

	query := `
		INSERT INTO user_registrations (
			user_registration_id,
			first_name,
			last_name,
			display_name,
			timezone,
			email
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) 
		RETURNING
			id,
			user_registration_id,
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
		user.FirstName,
		user.LastName,
		user.DisplayName,
		user.Timezone,
		user.Email,
	).Scan(
		&newUser.ID,
		&newUser.UserRegistrationID,
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

func UsersRegTableGetUserByID(userID string) (models.UserRegistration, error) {
	log.Println("- UsersRegTableGetUserByID")

	query := `
		SELECT 
			id,
			user_registration_id,
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

	return models.UserRegistration{}, nil
}

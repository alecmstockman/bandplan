package database

import (
	"bandplan/src/models"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func CreateUsersTable(db *sql.DB) error {
	log.Println("- CreateUsersTable")

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL UNIQUE,

		name TEXT NOT NULL,
		display_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		user_slug TEXT UNIQUE,

		password_hash TEXT NOT NULL,

		profile_image_id TEXT,
		profile_image_path TEXT,

		timezone TEXT,
		is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,

		last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Println("   Unable to create or load users table:", err)
		return err
	}

	return nil
}

func UsersTableCreateUser(name string, displayName string, email string, password string) (models.User, error) {
	log.Println("- UsersTableCreateUser")

	newID := uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("   Unable to hash password for user %s", name)
		return models.User{}, err
	}

	hashedPassword := string(hash)

	query := `
	INSERT INTO users (
		user_id,
		name,
		display_name,
		email,
		password_hash
	) VALUES (
		$1, $2, $3, $4, $5
	)
	RETURNING
		id,
		user_id,
		name,
		display_name,
		email,
		COALESCE(user_slug, ''),
		password_hash,
		COALESCE(profile_image_id, ''),
		COALESCE(profile_image_path, ''),
		COALESCE(timezone, ''),
		is_email_verified,
		last_login,
		created_at,
		updated_at
	`

	var newUser models.User

	err = DB.QueryRow(
		query,
		newID,
		name,
		displayName,
		email,
		hashedPassword,
	).Scan(
		&newUser.ID,
		&newUser.UserID,
		&newUser.Name,
		&newUser.DisplayName,
		&newUser.Email,
		&newUser.UserSlug,
		&newUser.PasswordHash,
		&newUser.ProfileImageID,
		&newUser.ProfileImagePath,
		&newUser.TimeZone,
		&newUser.IsEmailVerified,
		&newUser.LastLogin,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		log.Println("   UsersTableCreateUser err:", err)
		return models.User{}, err
	}

	return newUser, nil
}

func UsersTableGetUserByEmail(email string) (models.User, error) {
	log.Println("- UsersTableGetUserByEmail")
	log.Println("   email:", email)

	var user models.User

	query := `
	SELECT 
		id,
		user_id,
		name,
		display_name,
		email,
		COALESCE(user_slug, ''),
		password_hash,
		COALESCE(profile_image_id, ''),
		COALESCE(profile_image_path, ''),
		COALESCE(timezone, ''),
		is_email_verified,
		last_login,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	LIMIT 1
	`

	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.DisplayName,
		&user.Email,
		&user.UserSlug,
		&user.PasswordHash,
		&user.ProfileImageID,
		&user.ProfileImagePath,
		&user.TimeZone,
		&user.IsEmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("   UsersTableGetUserByEmail err:", err)
		return models.User{}, err
	}

	return user, nil
}

func UsersTableUpdateProfileImage(userID string, imageID string, imagePath string) error {
	log.Println("- UsersTableUpdateProfileImage")

	query := `
	UPDATE users
	SET
		profile_image_id = $1,
		profile_image_path = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE user_id = $3
	`

	_, err := DB.Exec(query, imageID, imagePath, userID)
	if err != nil {
		log.Println("   Unable to update user profile image:", err)
		return err
	}

	return nil
}

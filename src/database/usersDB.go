package database

import (
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,

		band_name TEXT,
		role TEXT NOT NULL DEFAULT 'Member',

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Unable to create or load users table")
		log.Fatal(err)
	}

	return nil
}

func UsersTableCreateUser(name string, bandName string, email string, password string) (models.User, error) {
	fmt.Println("\nUsersTableCreateUser")
	new_id := uuid.New()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Printf("Unable to save user %s", name)
		return models.User{}, err
	}

	hashedPassword := string(hash)

	query := `
	INSERT INTO users(
		user_id,
		name,
		email,
		password_hash,
		band_name
	) VALUES (
	 $1, $2, $3, $4, $5
	)
	RETURNING id, user_id, name, email, password_hash, band_name, role, created_at, updated_at
	`
	var newUser models.User

	err = DB.QueryRow(
		query,
		new_id,
		name,
		email,
		hashedPassword,
		bandName,
	).Scan(
		&newUser.ID,
		&newUser.UserID,
		&newUser.Name,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.BandName,
		&newUser.Role,
		&newUser.CreatedAt,
		&newUser.CreatedAt,
	)
	if err != nil {
		fmt.Println("UsersTableCreateUsers err: ", err)
		return models.User{}, err
	}

	fmt.Println(newUser)

	return newUser, nil
}

func UsersTableGetUserByEmail(email string) (models.User, error) {
	fmt.Printf("Getting user from users db: %s ", email)

	var user models.User

	query := `
	SELECT * 
	FROM users
	WHERE email = $1
	LIMIT 1
	`
	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.BandName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

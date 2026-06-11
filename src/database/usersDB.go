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
	fmt.Println("CreateUsersTable")
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL UNIQUE,
		band TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func UsersTableCreateUser(name string, band string, email string, password string) (models.User, error) {
	fmt.Println("\nUsersTableCreateUsers")
	id := uuid.New()

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
		band,
		email,
		password_hash
	) VALUES (
	 $1, $2, $3, $4, $5
	)
	RETURNING id, user_id, name, band, email, password_hash, created_at
	`
	var newUser models.User

	err = DB.QueryRow(
		query,
		id,
		name,
		band,
		email,
		hashedPassword,
	).Scan(
		&newUser.ID,
		&newUser.UserID,
		&newUser.Name,
		&newUser.Band,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.CreatedAt,
	)
	if err != nil {
		fmt.Println("nUsersTableCreateUsers err: ", err)
		return models.User{}, err
	}

	fmt.Println(newUser)

	return newUser, nil
}

func UsersTableGetUserByEmail(email string) (models.User, error) {
	fmt.Printf("Getting user from users db: %s", email)

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
		&user.Band,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

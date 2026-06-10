package database

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		band TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Id           int
	Name         string
	Band         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func UsersTableCreateUser(name string, band string, email string, password string) (User, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Printf("Unable to save user %s", name)
		return User{}, err
	}

	hashedPassword := string(hash)

	query := `
	INSERT INTO users(
		name,
		band,
		email,
		password_hash
	) VALUES ($1, $2, $3, $4)
	RETURNING id, name, band, email, password_hash, created_at
	`
	var newUser User

	err = DB.QueryRow(
		query,
		name,
		band,
		email,
		hashedPassword,
	).Scan(
		&newUser.Id,
		&newUser.Name,
		&newUser.Band,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Session struct {
	Id        int
	Email     string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func CreateSesssionsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXITS sessions (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		token TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL
		created_at DEFAULE CURRENT TIMESTAMP
	);
	`

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func SessionsTableInsertSession(email string, token string, expiresAt time.Time) (Session, error) {

	query := `
	INSERT INTO sessions (
		email
		token
		createdAt
	)
	VALUES (
		$1, $2, $3,
	) RETURNING id, email, token, expires_at, created_at
	`

	var session Session

	err := DB.QueryRow(
		query,
		email,
		token,
		expiresAt,
	).Scan(
		&session.Id,
		&session.Email,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func SessionsTableGetSessionByEmail(email string) (Session, error) {
	fmt.Println("SessionsTableGetTokenByEmail")

	var session Session

	query := `
	SELECT * 
	FROM sessions
	WHERE email = $1
	`

	err := DB.QueryRow(query, email).Scan(
		&session.Id,
		&session.Email,
		&session.Token,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		return Session{}, err
	}

	return Session{}, nil
}

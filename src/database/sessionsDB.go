package database

import (
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func CreateSesssionsTable(db *sql.DB) error {
	fmt.Println("CreateSessionsTable")

	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		users_id TEXT NOT NULL REFERENCES users(user_id),
		token TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func SessionsTableCreateSession(userID string, token string) (models.Session, error) {
	fmt.Println("\n SessionsTableCreateSession")
	fmt.Println("user_id: ", userID)
	fmt.Println("token: ", token)

	expires := time.Now()

	query := `
	INSERT INTO sessions (
		users_id,
		token,
		expires_at
	)
	VALUES (
		$1, $2, $3
	) RETURNING id, users_id, token, expires_at, created_at
	`

	var session models.Session

	err := DB.QueryRow(
		query,
		userID,
		token,
		expires,
	).Scan(
		&session.ID,
		&session.UsersID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		fmt.Println("create session err: ", err)
		return models.Session{}, err
	}
	return session, nil
}

func SessionsTableGetSessionByUserID(userID string) (models.Session, error) {
	fmt.Println("SessionsTableGetTokenByEmail")

	var session models.Session

	query := `
	SELECT * 
	FROM sessions
	WHERE user_id = $1
	`

	err := DB.QueryRow(query, userID).Scan(
		&session.ID,
		&session.UsersID,
		&session.Token,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func SessionsTableGetValidateToken(token string) (bool, error) {
	var validated bool

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM sessions
		WHERE token = $1
		AND expires_at > NOW()
	)
	`
	err := DB.QueryRow(query, token).Scan(&validated)

	if err != nil {
		return false, err
	}

	return validated, nil
}

func SessionsTableGetSessionByToken(token string) (models.Session, error) {
	var session models.Session

	query := `
	SELECT * 
	FROM sessions
	WHERE token = $1
	`
	err := DB.QueryRow(
		query, token,
	).Scan(
		&session.ID,
		&session.UsersID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func SessionsTableGetUserByToken(token string) (models.User, error) {
	var user models.User

	query := `
	SELECT
		users.id,
		users.user_id,
		users.name,
		users.email,
		users.password_hash,
		users.created_at
	FROM users
	LEFT JOIN sessions
	ON users.user_id = sessions.users_id
	WHERE token = $1 
	`
	err := DB.QueryRow(
		query, token,
	).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func SessionsTableDeleteSessionByToken(token string) error {
	query := `
	DELETE FROM sessions
	WEHRE token = $1
	`
	_, err := DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

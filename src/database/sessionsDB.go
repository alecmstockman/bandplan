package database

import (
	"bandplan/src/models"
	"database/sql"
	"log"
	"time"
)

func CreateSesssionsTable(db *sql.DB) error {
	log.Println("- CreateSesssionsTable")
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		users_id TEXT NOT NULL REFERENCES users(user_id),
		band_id TEXT REFERENCES bands(band_id),
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Println("   Unable to create or load sessions table")
		log.Fatal(err)
	}
	return nil
}

func SessionsTableCreateSession(c models.CreateSessionParams) (models.Session, error) {
	log.Println("- SessionsTableCreateSession")

	expires := time.Now()

	query := `
	INSERT INTO sessions (
		users_id,
		band_id,
		token,
		expires_at
	)
	VALUES (
		$1, $2, $3, $4
	) RETURNING id, users_id, band_id, token, created_at, expires_at
	`

	var session models.Session

	err := DB.QueryRow(
		query,
		c.UserID,
		c.BandID,
		c.Token,
		expires,
	).Scan(
		&session.ID,
		&session.UsersID,
		&session.BandID,
		&session.Token,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		log.Println("   create session err: ", err)
		return models.Session{}, err
	}
	return session, nil
}

func SessionsTableGetSessionByUserID(userID string) (models.Session, error) {
	log.Println("- SessionsTableGetSessionByUserID")

	var session models.Session

	query := `
	SELECT 
		id,
		users_id,
		COALESCE(band_id, ''),
		token,
		created_at,
		expires_at
	FROM sessions
	WHERE user_id = $1
	`

	err := DB.QueryRow(query, userID).Scan(
		&session.ID,
		&session.UsersID,
		&session.BandID,
		&session.Token,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		log.Println("   Unable to get session by user id: ", err)
		return models.Session{}, err
	}

	return session, nil
}

func SessionsTableGetValidatedBYToken(token string) (bool, error) {
	log.Println("- SessionsTableGetValidatedBYToken")
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
		log.Println("   Unable to get valid token by token: ", err)
		return false, err
	}

	return validated, nil
}

func SessionsTableGetSessionByToken(token string) (models.Session, error) {
	log.Println("- SessionsTableGetSessionByToken")
	var session models.Session

	query := `
	SELECT
		id,
		users_id,
		COALESCE(band_id, ''),
		token,
		created_at,
		expires_at
	FROM sessions
	WHERE token = $1
	`
	err := DB.QueryRow(
		query, token,
	).Scan(
		&session.ID,
		&session.UsersID,
		&session.Token,
		&session.BandID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		log.Println("   Unable to get session by token from sessions db: ", err)
		return models.Session{}, err
	}

	return session, nil
}

func SessionsTableGetUserByToken(token string) (models.User, error) {
	log.Println("- SessionsTableGetUserByToken")
	var user models.User

	query := `
	SELECT
		users.id,
		users.user_id,
		users.name,
		users.display_name,
		users.email,
		users.user_slug,
		users.password_hash,
		users.is_admin,
		COALESCE(users.profile_image_id, ''),
		COALESCE(users.profile_image_path, ''),
		COALESCE(users.timezone, ''), 
		users.is_email_verified,
		users.last_login,
		users.created_at,
		users.updated_at
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
		&user.DisplayName,
		&user.Email,
		&user.UserSlug,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.ProfileImageID,
		&user.ProfileImagePath,
		&user.TimeZone,
		&user.IsEmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("   Error authenticating user: ", err)
		return models.User{}, err
	}

	return user, nil
}

func SessionsTableDeleteSessionByToken(token string) error {
	log.Println("- SessionsTableDeleteSessionByToken")
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

func SessionsTableGetAuthContextByToken(token string) (models.User, models.Band, error) {
	log.Println("- SessionsTableGetAuthContextByToken")

	query := `
		SELECT
			u.id,
			u.user_id,
			u.name,
			u.display_name,
			u.email,
			u.is_admin,
			u.profile_image_id,
			u.profile_image_path,
			u.timezone,
			u.is_email_verified,
			u.last_login,
			u.created_at,
			u.updated_at,

			b.id,
			b.band_id,
			b.name,
			b.created_at

		fROM sessions s

		JOIN users u
			ON u.user_id = s.users_id
		LEFT JOIN bands b
			ON b.band_id = s.band_id

		WHERE s.token = $1
		AND s.expires_at > CURRENT_TIMESTAMP
	`

	var user models.User
	var band models.Band

	err := DB.QueryRow(query, token).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.DisplayName,
		&user.Email,
		&user.IsAdmin,
		&user.ProfileImageID,
		&user.ProfileImagePath,
		&user.TimeZone,
		&user.IsEmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,

		&band.ID,
		&band.BandID,
		&band.Name,
		&band.CreatedAt,
	)
	if err != nil {
		log.Println("   Unable to get user and band from db: ", err)
		return models.User{}, models.Band{}, err
	}
	return user, band, nil
}

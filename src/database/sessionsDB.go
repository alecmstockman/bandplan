package database

import (
	"bandplan/src/models"
	"errors"
	"time"
)

func SessionsTableCreateSession(c models.CreateSessionParams) (models.Session, error) {

	expires := time.Now().Add(1 * time.Hour)

	query := `
	INSERT INTO sessions (
		user_id,
		band_id,
		token,
		expires_at
	)
	VALUES (
		$1, $2, $3, $4
	) RETURNING id, user_id, band_id, token, created_at, expires_at
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
		return models.Session{}, err
	}
	return session, nil
}

func SessionsTableGetValidatedBYToken(token string) (bool, error) {

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
		SELECT
			id,
			user_id,
			COALESCE(band_id, ''),
			token,
			created_at,
			expires_at
		FROM sessions
		WHERE token = $1
			AND expires_at > NOW()
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
		users.display_name,
		users.email,
		users.slug,
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
	LEFT JOIN sessions s
	ON users.user_id = s.user_id
	WHERE token = $1 
		AND s.expires_at > NOW()
	`
	err := DB.QueryRow(
		query, token,
	).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.DisplayName,
		&user.Email,
		&user.Slug,
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
		return models.User{}, err
	}

	return user, nil
}

func SessionsTableDeleteSessionByUserID(userID, token string) (bool, error) {

	query := `
		DELETE FROM sessions
		WHERE user_id = $1
			AND token = $2
			AND expires_at > NOW()
	`
	affected, err := DB.Exec(query, userID, token)
	if err != nil {
		return false, err
	}

	result, err := affected.RowsAffected()
	if err != nil {
		return false, err
	}

	if result != 1 {
		return false, errors.New("Unable to delete session by token and userID")
	}

	return true, nil
}

func SessionsTableDeleteSessionByToken(token string) error {

	query := `
		DELETE FROM sessions
		WHERE token = $1
	`
	_, err := DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

func SessionsTableGetAuthContextByToken(token string) (models.User, models.Band, error) {

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
			ON u.user_id = s.user_id
		LEFT JOIN bands b
			ON b.band_id = s.band_id

		WHERE s.token = $1
		AND s.expires_at > NOW()
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
		return models.User{}, models.Band{}, err
	}
	return user, band, nil
}

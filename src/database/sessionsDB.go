package database

import (
	"bandplan/src/models"
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

func SessionsTableDeleteSessionByToken(token string) error {

	query := `
		DELETE FROM sessions
		WHERE token = $1
	`
	result, err := DB.Exec(query, token)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return nil
	}
	return nil
}

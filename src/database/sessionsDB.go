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
	) RETURNING id, user_id, band_id, created_at, expires_at
	`

	var session models.Session

	err := DB.QueryRow(
		query,
		c.UserID,
		c.BandID,
		c.TokenHash,
		expires,
	).Scan(
		&session.ID,
		&session.UsersID,
		&session.BandID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	session.Token = c.Token

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

func SessionsTableGetCurrentBandByToken(token, userID string) (models.Band, error) {
	query := `
	SELECT
		bands.id,
		bands.band_id,
		bands.name,
		bands.slug,
		bands.created_at,
		bands.created_by,
		bands.updated_at,
		bands.updated_by
	FROM sessions
	JOIN bands
		ON bands.band_id = sessions.band_id
	JOIN band_members
		ON band_members.band_id = bands.band_id
		AND band_members.user_id = sessions.user_id
	WHERE sessions.token = $1
		AND sessions.user_id = $2
		AND sessions.expires_at > NOW()
	`

	var band models.Band
	err := DB.QueryRow(query, token, userID).Scan(
		&band.ID,
		&band.BandID,
		&band.Name,
		&band.Slug,
		&band.CreatedAt,
		&band.CreatedBy,
		&band.UpdatedAt,
		&band.UpdatedBy,
	)
	if err != nil {
		return models.Band{}, err
	}

	return band, nil
}

func SessionsTableSetCurrentBand(token, userID, bandID string) (bool, error) {
	query := `
	UPDATE sessions
	SET band_id = $1
	WHERE token = $2
		AND user_id = $3
		AND expires_at > NOW()
		AND EXISTS (
			SELECT 1
			FROM band_members
			WHERE band_members.band_id = $1
				AND band_members.user_id = sessions.user_id
		)
	`

	result, err := DB.Exec(query, bandID, token, userID)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
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

package database

import (
	"bandplan/src/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func SessionsTableCreateSession(c models.CreateSessionParams) (models.Session, error) {

	expires := time.Now().Add(30 * 24 * time.Hour)

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
		return models.Session{}, fmt.Errorf("SessionsTableCreateSession - create session:  %w", err)
	}
	return session, nil
}

func SessionsTableGetUserByToken(token string) (models.User, error) {

	var user models.User

	query := `
	SELECT
		u.id,
		u.user_id,
		COALESCE(u.name, ''),
		COALESCE(u.first_name, ''),
		COALESCE(u.last_name, ''),
		u.display_name,
		u.email,
		COALESCE(u.slug, ''),
		u.password_hash,
		COALESCE(bm.is_admin, FALSE),
		COALESCE(bm.edit_band, FALSE),
		COALESCE(bm.manage_members, FALSE),
		COALESCE(bm.create_songs, FALSE),
		COALESCE(bm.edit_songs, FALSE),
		COALESCE(bm.delete_songs, FALSE),
		COALESCE(bm.create_setlists, FALSE),
		COALESCE(bm.edit_setlists, FALSE),
		COALESCE(bm.delete_setlists, FALSE),
		COALESCE(bm.create_events, FALSE),
		COALESCE(bm.edit_events, FALSE),
		COALESCE(bm.delete_events, FALSE),
		COALESCE(bm.create_todos, FALSE),
		COALESCE(bm.edit_todos, FALSE),
		COALESCE(bm.delete_todos, FALSE),
		COALESCE(u.profile_image_id, ''),
		COALESCE(u.profile_image_path, ''),
		COALESCE(u.timezone, ''), 
		u.is_email_verified,
		COALESCE(u.last_login, NOW()),
		u.created_at,
		u.updated_at
	FROM users u
	LEFT JOIN sessions s
		ON u.user_id = s.user_id
	LEFT JOIN band_members bm
		ON u.user_id = bm.user_id
		AND s.band_id = bm.band_id
	WHERE token = $1 
		AND s.expires_at > NOW()
	`
	err := DB.QueryRow(
		query, token,
	).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.FirstName,
		&user.LastName,
		&user.DisplayName,
		&user.Email,
		&user.Slug,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.Permissions.EditBand,
		&user.Permissions.ManageMembers,
		&user.Permissions.CreateSongs,
		&user.Permissions.EditSongs,
		&user.Permissions.DeleteSongs,
		&user.Permissions.CreateSetlists,
		&user.Permissions.EditSetlists,
		&user.Permissions.DeleteSetlists,
		&user.Permissions.CreateEvents,
		&user.Permissions.EditEvents,
		&user.Permissions.DeleteEvents,
		&user.Permissions.CreateTodos,
		&user.Permissions.EditTodos,
		&user.Permissions.DeleteTodos,
		&user.ProfileImageID,
		&user.ProfileImagePath,
		&user.TimeZone,
		&user.IsEmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, fmt.Errorf("SessionsTableGetUserByToken - get user by token:  %w", err)
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
			COALESCE(bands.profile_image_id, ''),
			COALESCE(bands.profile_image_path, ''),
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
		&band.ProfileImageID,
		&band.ProfileImagePath,
		&band.CreatedAt,
		&band.CreatedBy,
		&band.UpdatedAt,
		&band.UpdatedBy,
	)
	if err != nil {
		return models.Band{}, fmt.Errorf("SessionsTableGetCurrentBandByToken - get current band by token:  %w", err)
	}

	return band, nil
}

func SessionsTableGetUserAndBandByToken(ctx context.Context, token string) (models.User, models.Band, error) {
	var user models.User
	var bandID sql.NullString
	var bandDatabaseID sql.NullInt64
	var bandName sql.NullString
	var bandSlug sql.NullString
	var bandProfileImageID sql.NullString
	var bandProfileImagePath sql.NullString
	var bandCreatedAt sql.NullTime
	var bandCreatedBy sql.NullString
	var bandUpdatedAt sql.NullTime
	var bandUpdatedBy sql.NullString

	query := `
		SELECT
			u.id,
			u.user_id,
			COALESCE(u.name, ''),
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			u.display_name,
			u.email,
			COALESCE(u.slug, ''),
			u.password_hash,
			COALESCE(bm.is_admin, FALSE),
			COALESCE(bm.edit_band, FALSE),
			COALESCE(bm.manage_members, FALSE),
			COALESCE(bm.create_songs, FALSE),
			COALESCE(bm.edit_songs, FALSE),
			COALESCE(bm.delete_songs, FALSE),
			COALESCE(bm.create_setlists, FALSE),
			COALESCE(bm.edit_setlists, FALSE),
			COALESCE(bm.delete_setlists, FALSE),
			COALESCE(bm.create_events, FALSE),
			COALESCE(bm.edit_events, FALSE),
			COALESCE(bm.delete_events, FALSE),
			COALESCE(bm.create_todos, FALSE),
			COALESCE(bm.edit_todos, FALSE),
			COALESCE(bm.delete_todos, FALSE),
			COALESCE(u.profile_image_id, ''),
			COALESCE(u.profile_image_path, ''),
			COALESCE(u.timezone, ''), 
			u.is_email_verified,
			COALESCE(u.last_login, NOW()),
			u.created_at,
			u.updated_at,
			b.id,
			b.band_id,
			b.name,
			b.slug,
			b.profile_image_id,
			b.profile_image_path,
			b.created_at,
			b.created_by,
			b.updated_at,
			b.updated_by
		FROM sessions s
		JOIN users u
			ON u.user_id = s.user_id
		LEFT JOIN band_members bm
			ON bm.user_id = s.user_id
			AND bm.band_id = s.band_id
		LEFT JOIN bands b
			ON b.band_id = bm.band_id
		WHERE s.token = $1
			AND s.expires_at > NOW()
	`

	err := DB.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.FirstName,
		&user.LastName,
		&user.DisplayName,
		&user.Email,
		&user.Slug,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.Permissions.EditBand,
		&user.Permissions.ManageMembers,
		&user.Permissions.CreateSongs,
		&user.Permissions.EditSongs,
		&user.Permissions.DeleteSongs,
		&user.Permissions.CreateSetlists,
		&user.Permissions.EditSetlists,
		&user.Permissions.DeleteSetlists,
		&user.Permissions.CreateEvents,
		&user.Permissions.EditEvents,
		&user.Permissions.DeleteEvents,
		&user.Permissions.CreateTodos,
		&user.Permissions.EditTodos,
		&user.Permissions.DeleteTodos,
		&user.ProfileImageID,
		&user.ProfileImagePath,
		&user.TimeZone,
		&user.IsEmailVerified,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&bandDatabaseID,
		&bandID,
		&bandName,
		&bandSlug,
		&bandProfileImageID,
		&bandProfileImagePath,
		&bandCreatedAt,
		&bandCreatedBy,
		&bandUpdatedAt,
		&bandUpdatedBy,
	)
	if err != nil {
		return models.User{}, models.Band{}, fmt.Errorf("SessionsTableGetUserAndBandByToken - get user and band by token:  %w", err)
	}

	if !bandID.Valid {
		return user, models.Band{}, fmt.Errorf(
			"SessionsTableGetUserAndBandByToken - get current band by token:  %w",
			sql.ErrNoRows,
		)
	}

	band := models.Band{
		ID:               int(bandDatabaseID.Int64),
		BandID:           bandID.String,
		Name:             bandName.String,
		Slug:             bandSlug.String,
		ProfileImageID:   bandProfileImageID.String,
		ProfileImagePath: bandProfileImagePath.String,
		CreatedAt:        bandCreatedAt.Time,
		CreatedBy:        bandCreatedBy.String,
		UpdatedAt:        bandUpdatedAt.Time,
		UpdatedBy:        bandUpdatedBy.String,
	}

	return user, band, nil
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
		return false, fmt.Errorf("SessionsTableSetCurrentBand - set current band:  %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("SessionsTableSetCurrentBand - count updated sessions:  %w", err)
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
		return fmt.Errorf("SessionsTableDeleteSessionByToken - delete session by token:  %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SessionsTableDeleteSessionByToken - count deleted sessions:  %w", err)
	}

	if affected == 0 {
		return nil
	}
	return nil
}

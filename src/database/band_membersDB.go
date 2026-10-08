package database

import (
	"bandplan/src/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrBandAlreadyJoined = errors.New("user is already a member of this band")

func BandMembersCreateMember(bandID string, userID string) error {

	query := `
	INSERT INTO band_members(
		band_id,
		user_id
	) VALUES ($1, $2)
	`
	_, err := DB.Exec(
		query,
		bandID,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func BandMembersJoinWithAccessCode(ctx context.Context, userID, sessionToken, accessCodeHash string) (string, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var bandID string
	err = tx.QueryRowContext(ctx, `
		DELETE FROM access_codes
		WHERE code_hash = $1
			AND expires_at > NOW()
		RETURNING band_id
	`, accessCodeHash).Scan(&bandID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrAccessCodeExpired
	}
	if err != nil {
		return "", fmt.Errorf("consume access code: %w", err)
	}

	result, err := tx.ExecContext(ctx, `
		INSERT INTO band_members (band_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (band_id, user_id) DO NOTHING
	`, bandID, userID)
	if err != nil {
		return "", fmt.Errorf("insert band member: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("count inserted band members: %w", err)
	}
	if affected == 0 {
		return "", ErrBandAlreadyJoined
	}

	var chatID string
	err = tx.QueryRowContext(ctx, `
		SELECT chat_id
		FROM chats
		WHERE band_id = $1
			AND is_primary = TRUE
	`, bandID).Scan(&chatID)
	if err != nil {
		return "", fmt.Errorf("get primary chat: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`, chatID, userID)
	if err != nil {
		return "", fmt.Errorf("insert primary chat member: %w", err)
	}

	result, err = tx.ExecContext(ctx, `
		UPDATE sessions
		SET band_id = $1
		WHERE token = $2
			AND user_id = $3
			AND expires_at > NOW()
	`, bandID, sessionToken, userID)
	if err != nil {
		return "", fmt.Errorf("update current band: %w", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("count updated sessions: %w", err)
	}
	if affected != 1 {
		return "", errors.New("active session not found")
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return bandID, nil
}

func BandMembersGetMembersByBandID(bandID string) ([]models.User, error) {

	query := `
	SELECT
		u.id,
		u.user_id,
		u.name,
		u.display_name,
		u.email,
		u.slug,
		u.password_hash,
		u.is_admin,
		COALESCE(u.profile_image_id, ''),
		COALESCE(u.profile_image_path, ''),
		COALESCE(u.timezone, ''),
		u.is_email_verified,
		u.last_login,
		u.created_at,
		u.updated_at

		FROM band_members b
		JOIN users u
			ON u.user_id = b.user_id
		WHERE b.band_id = $1
	`

	rows, err := DB.Query(query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
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
			return []models.User{}, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func BandMembersGetMemberNameAndID(userID, bandID string) ([]models.User, error) {
	query := `
		SELECT
			u.user_id,
			u.display_name
		FROM band_members bm
		JOIN users u
			ON u.user_id = bm.user_id
		WHERE bm.band_id = $1
			AND EXISTS (
				SELECT 1
				FROM band_members requester
				WHERE requester.band_id = bm.band_id
					AND requester.user_id = $2
			)
		ORDER BY u.display_name
	`

	rows, err := DB.Query(query, bandID, userID)
	if err != nil {
		return nil, fmt.Errorf("query band member names and IDs: %w", err)
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var member models.User
		if err := rows.Scan(&member.UserID, &member.DisplayName); err != nil {
			return nil, fmt.Errorf("scan band member name and ID: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate band member names and IDs: %w", err)
	}

	return members, nil
}

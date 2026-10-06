package database

import (
	"bandplan/src/models"
	"fmt"
)

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

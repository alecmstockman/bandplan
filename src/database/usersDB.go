package database

import (
	"bandplan/src/models"
	"strings"

	_ "github.com/lib/pq"
)

func UsersTableGetUserByEmail(email string) (models.User, error) {
	var user models.User

	query := `
	SELECT 
		id,
		user_id,
		name,
		display_name,
		LOWER(email),
		COALESCE(slug, ''),
		password_hash,
		is_admin,
		COALESCE(profile_image_id, ''),
		COALESCE(profile_image_path, ''),
		COALESCE(timezone, ''),
		is_email_verified,
		legal_accepted,
		legal_accepted_at,
		last_login,
		created_at,
		updated_at
	FROM users
	WHERE LOWER(email) = $1
	LIMIT 1
	`

	err := DB.QueryRow(
		query,
		strings.ToLower(email),
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
		&user.LegalAccepted,
		&user.LegalAcceptedAt,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UsersTableGetUserByUserID(userID string) (models.User, error) {
	var user models.User

	query := `
	SELECT
		id,
		user_id,
		name,
		first_name,
		COALESCE(last_name, ''),
		display_name,
		email,
		COALESCE(slug, ''),
		password_hash,
		is_admin,
		COALESCE(profile_image_id, ''),
		COALESCE(profile_image_path, ''),
		COALESCE(timezone, ''),
		is_email_verified,
		legal_accepted,
		legal_accepted_at,
		last_login,
		created_at,
		updated_at
	FROM users
	WHERE user_id = $1
	`

	err := DB.QueryRow(query, userID).Scan(
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
		&user.ProfileImageID,
		&user.ProfileImagePath,
		&user.TimeZone,
		&user.IsEmailVerified,
		&user.LegalAccepted,
		&user.LegalAcceptedAt,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UsersTableUpdateProfileImage(userID string, imageID string, imagePath string) error {
	query := `
	UPDATE users
	SET
		profile_image_id = $1,
		profile_image_path = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE user_id = $3
	`

	_, err := DB.Exec(query, imageID, imagePath, userID)
	if err != nil {
		return err
	}

	return nil
}

func UsersTableGetUsersByBand(bandID string) ([]models.User, error) {

	query := `
		SELECT
			u.id,
			u.user_id,
			u.name,
			u.display_name,
			u.email,
			COALESCE(u.slug, ''),
			u.password_hash,
			u.is_admin,
			COALESCE(u.profile_image_id, ''),
			COALESCE(u.profile_image_path, ''),
			COALESCE(u.timezone, ''),
			u.is_email_verified,
			u.legal_accepted,
			u.legal_accepted_at,
			u.last_login,
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN band_members bm
			ON u.user_id = bm.user_id
		WHERE bm.band_id = $1
	`

	rows, err := DB.Query(query, bandID)
	if err != nil {
		return []models.User{}, err
	}

	defer rows.Close()

	var members []models.User

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
			&user.LegalAccepted,
			&user.LegalAcceptedAt,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return []models.User{}, err
		}

		members = append(members, user)
	}
	return members, nil
}

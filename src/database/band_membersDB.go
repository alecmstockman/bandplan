package database

import (
	"bandplan/src/models"
	"log"
)

func BandMembersCreateMember(bandID string, userID string) error {
	log.Println("- BandMembersCreateMember")

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
		log.Println("   Unable to create band member: ", err)
		return err
	}

	return nil
}

func BandMembersGetMembersByBandID(bandID string) ([]models.User, error) {
	log.Println("- BandMembersGetMembersByBandID")

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
		log.Println("   Unable to query band_members from database: ", err)
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
			log.Println("   Unable to get band members: ", err)
			return []models.User{}, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

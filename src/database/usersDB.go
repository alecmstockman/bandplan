package database

import (
	"bandplan/src/models"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func UsersTableCreateUser(name string, displayName string, slug string, email string, password string, isAdmin bool) (models.User, error) {
	log.Println("- UsersTableCreateUser")

	newID := uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("   Unable to hash password for user %s", name)
		return models.User{}, err
	}

	hashedPassword := string(hash)

	query := `
	INSERT INTO users (
		user_id,
		name,
		display_name,
		email,
		slug,
		password_hash,
		is_admin
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	)
	RETURNING
		id,
		user_id,
		name,
		display_name,
		email,
		slug,
		password_hash,
		is_admin,
		COALESCE(profile_image_id, ''),
		COALESCE(profile_image_path, ''),
		COALESCE(timezone, ''),
		is_email_verified,
		last_login,
		created_at,
		updated_at
	`

	var newUser models.User

	err = DB.QueryRow(
		query,
		newID,
		name,
		displayName,
		email,
		slug,
		hashedPassword,
		isAdmin,
	).Scan(
		&newUser.ID,
		&newUser.UserID,
		&newUser.Name,
		&newUser.DisplayName,
		&newUser.Email,
		&newUser.Slug,
		&newUser.PasswordHash,
		&newUser.IsAdmin,
		&newUser.ProfileImageID,
		&newUser.ProfileImagePath,
		&newUser.TimeZone,
		&newUser.IsEmailVerified,
		&newUser.LastLogin,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		log.Println("   Unable to create user in table:", err)
		return models.User{}, err
	}

	return newUser, nil
}

func UsersTableGetUserByEmail(email string) (models.User, error) {
	log.Println("- UsersTableGetUserByEmail")
	log.Println("   email:", email)

	var user models.User

	query := `
	SELECT 
		id,
		user_id,
		name,
		display_name,
		email,
		COALESCE(slug, ''),
		password_hash,
		is_admin,
		COALESCE(profile_image_id, ''),
		COALESCE(profile_image_path, ''),
		COALESCE(timezone, ''),
		is_email_verified,
		last_login,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	LIMIT 1
	`

	err := DB.QueryRow(query, email).Scan(
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
		log.Println("   UsersTableGetUserByEmail err:", err)
		return models.User{}, err
	}

	return user, nil
}

func UsersTableUpdateProfileImage(userID string, imageID string, imagePath string) error {
	log.Println("- UsersTableUpdateProfileImage")

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
		log.Println("   Unable to update user profile image:", err)
		return err
	}

	return nil
}

func UsersTableGetUsersByBand(bandID string) ([]models.User, error) {
	log.Println("- UsersTableGetUsersByBand")

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
		log.Println("   Unable to get users band bandID: ", err)
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
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("   Unable to get users from database: ", err)
			return []models.User{}, err
		}

		members = append(members, user)
	}
	return members, nil
}

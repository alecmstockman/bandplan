package database

import (
	"bandplan/src/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func BandsTableCreateBand(bandName string, userID string, bandSlug string) (models.Band, error) {

	newBandID := uuid.New().String()

	query := `
	INSERT INTO bands(
		band_id,
		name,
		slug,
		created_by,
		updated_by
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING id, band_id, name, slug, created_at, created_by, updated_at, updated_by
	`
	var newBand models.Band

	err := DB.QueryRow(
		query,
		newBandID,
		bandName,
		bandSlug,
		userID,
		userID,
	).Scan(
		&newBand.ID,
		&newBand.BandID,
		&newBand.Name,
		&newBand.Slug,
		&newBand.CreatedAt,
		&newBand.CreatedBy,
		&newBand.UpdatedAt,
		&newBand.UpdatedBy,
	)
	if err != nil {
		return models.Band{}, fmt.Errorf("BandsTableCreateBand - create band:  %w", err)
	}
	return newBand, nil
}

func BandsTableGetBandByName(bandName string) (models.Band, error) {

	query := `
	SELECT
		id,
		band_id,
		name,
		slug,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM bands
	WHERE name = $1
	LIMIT 1
	`
	var band models.Band

	err := DB.QueryRow(query, bandName).Scan(
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
		return models.Band{}, fmt.Errorf("BandsTableGetBandByName - get band by name:  %w", err)
	}

	return band, nil
}

func BandsTableGetBandByBandID(bandID string) (models.Band, error) {

	query := `
	SELECT
		id,
		band_id,
		name,
		slug,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM bands
	WHERE band_id = $1
	LIMIT 1
	`
	var band models.Band

	err := DB.QueryRow(query, bandID).Scan(
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
		return models.Band{}, fmt.Errorf("BandsTableGetBandByBandID - get band by band id:  %w", err)
	}

	return band, nil
}

func BandsTableGetBandByUserID(userID string) (models.Band, error) {

	var band models.Band

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
	FROM bands
	JOIN band_members 
		ON bands.band_id = band_members.band_id
	WHERE band_members.user_id = $1
	ORDER BY band_members.created_at, bands.band_id
	LIMIT 1
	`
	err := DB.QueryRow(query, userID).Scan(
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
		return models.Band{}, fmt.Errorf("BandsTableGetBandByUserID - get band by user id:  %w", err)
	}
	return band, nil
}

func BandsTableGetBandsByUserID(userID string) ([]models.Band, error) {
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
	FROM bands
	JOIN band_members
		ON bands.band_id = band_members.band_id
	WHERE band_members.user_id = $1
	ORDER BY bands.name, bands.band_id
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("BandsTableGetBandsByUserID - query bands by user id:  %w", err)
	}
	defer rows.Close()

	bands := make([]models.Band, 0)
	for rows.Next() {
		var band models.Band
		err = rows.Scan(
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
			return nil, fmt.Errorf("BandsTableGetBandsByUserID - scan band:  %w", err)
		}
		bands = append(bands, band)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("BandsTableGetBandsByUserID - iterate bands:  %w", err)
	}

	return bands, nil
}

func BandsTableGetBandByBandIDAndUserID(bandID, userID string) (models.Band, error) {
	query := `
		SELECT
			bands.id,
			bands.band_id,
			bands.name,
			bands.slug,
			COALESCE(bands.profile_image_id, ''),
			COALESCE(bands.profile_image_path, ''),
			COALESCE(bands.email_address, ''),
			bands.is_email_verified,
			COALESCE(bands.timezone, ''),
			COALESCE(bands.city, ''),
			COALESCE(bands.state, ''),
			COALESCE(bands.country, ''),
			bands.created_at,
			bands.created_by,
			bands.updated_at,
			bands.updated_by
		FROM bands
		JOIN band_members
			ON bands.band_id = band_members.band_id
		WHERE bands.band_id = $1
			AND band_members.user_id = $2
	`

	var band models.Band
	err := DB.QueryRow(query, bandID, userID).Scan(
		&band.ID,
		&band.BandID,
		&band.Name,
		&band.Slug,
		&band.ProfileImageID,
		&band.ProfileImagePath,
		&band.EmailAddress,
		&band.IsEmailVerified,
		&band.Timezone,
		&band.City,
		&band.State,
		&band.Country,
		&band.CreatedAt,
		&band.CreatedBy,
		&band.UpdatedAt,
		&band.UpdatedBy,
	)
	if err != nil {
		return models.Band{}, fmt.Errorf("BandsTableGetBandByBandIDAndUserID - get band by band id and user id:  %w", err)
	}

	return band, nil
}

func BandsTableUpdateProfileImage(bandID, imageID, imagePath string) error {
	query := `
	UPDATE bands
	SET
		profile_image_id = $1,
		profile_image_path = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE band_id = $3
	`

	_, err := DB.Exec(query, imageID, imagePath, bandID)
	if err != nil {
		return fmt.Errorf("BandsTableUpdateProfileImage - update profile image: %w", err)
	}

	return nil
}

func BandsTableUpdateProfile(ctx context.Context, band models.Band, userID string) (bool, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("BandsTableUpdateProfile - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentEmail string
	err = tx.QueryRowContext(ctx, `
		SELECT COALESCE(b.email_address, '')
		FROM bands b
		JOIN band_members bm ON bm.band_id = b.band_id
		WHERE b.band_id = $1
			AND bm.user_id = $2
			AND bm.is_admin = TRUE
		FOR UPDATE OF b
	`, band.BandID, userID).Scan(&currentEmail)
	if err != nil {
		return false, fmt.Errorf("BandsTableUpdateProfile - get current email: %w", err)
	}

	emailChanged := currentEmail != band.EmailAddress
	result, err := tx.ExecContext(ctx, `
		UPDATE bands
		SET
			name = $1,
			email_address = NULLIF($2, ''),
			timezone = $3,
			city = NULLIF($4, ''),
			state = NULLIF($5, ''),
			country = NULLIF($6, ''),
			is_email_verified = CASE WHEN $7 THEN FALSE ELSE is_email_verified END,
			email_verified_at = CASE WHEN $7 THEN NULL ELSE email_verified_at END,
			updated_by = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE band_id = $9
	`,
		band.Name,
		band.EmailAddress,
		band.Timezone,
		band.City,
		band.State,
		band.Country,
		emailChanged,
		userID,
		band.BandID,
	)
	if err != nil {
		return false, fmt.Errorf("BandsTableUpdateProfile - update profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("BandsTableUpdateProfile - count updated bands: %w", err)
	}
	if rowsAffected == 0 {
		return false, sql.ErrNoRows
	}

	if emailChanged {
		if _, err := tx.ExecContext(ctx, `DELETE FROM band_email_verification_tokens WHERE band_id = $1`, band.BandID); err != nil {
			return false, fmt.Errorf("BandsTableUpdateProfile - delete email verification token: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("BandsTableUpdateProfile - commit transaction: %w", err)
	}

	return emailChanged, nil
}

func BandsTableGetBandNameByID(bandID string) (string, error) {

	bandName := ""

	query := `
	SELECT name
	FROM bands
	WHERE band_id = $1
	`

	err := DB.QueryRow(query, bandID).Scan(&bandName)
	if err != nil {
		return "", fmt.Errorf("BandsTableGetBandNameByID - get band name by id:  %w", err)
	}

	return bandName, nil
}

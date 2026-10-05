package database

import (
	"bandplan/src/models"

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
		return models.Band{}, err
	}
	return newBand, nil
}

func BandsTableGetBandByName(bandName string) (models.Band, error) {

	query := `
	SELECT *
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
		return models.Band{}, err
	}

	return band, nil
}

func BandsTableGetBandByBandID(bandID string) (models.Band, error) {

	query := `
	SELECT *
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
		return models.Band{}, err
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
		bands.created_at,
		bands.created_by,
		bands.updated_at,
		bands.updated_by
	FROM bands
	JOIN band_members 
		ON bands.band_id = band_members.band_id
	WHERE band_members.user_id = $1
	LIMIT 1
	`
	err := DB.QueryRow(query, userID).Scan(
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

func BandsTableGetBandNameByID(bandID string) (string, error) {

	bandName := ""

	query := `
	SELECT name
	FROM bands
	WHERE band_id = $1
	`

	err := DB.QueryRow(query, bandID).Scan(&bandName)
	if err != nil {
		return "", err
	}

	return bandName, nil
}

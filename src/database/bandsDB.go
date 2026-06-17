package database

import (
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"log"
)

func CreateBandsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS bands (
		id SERIAL PRIMARY KEY,
		band_id TEXT NOT NULL UNIQUE,
		band_name TEXT NOT NULL, 
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Unable to create or load bands table")
		log.Fatal(err)
	}
	return nil
}

func BandsTableGetBandByName(bandName string) (models.Band, error) {
	fmt.Println("BandsTableGetBandByName")

	query := `
	SELECT *
	FROM bands
	WHERE band_name = $1
	LIMIT 1
	`

	var band models.Band

	err := DB.QueryRow(query, bandName).Scan(
		&band.ID,
		&band.BandID,
		&band.Name,
		&band.CreatedAt,
	)

	if err != nil {
		return models.Band{}, err
	}

	return band, nil
}

func BandsTableCreateBand(bandName string, userID string) (models.Band, error) {
	fmt.Println("BandsTableCreateBand")

	query := `
	INSERT INTO bands(
		bandName
	) VALUES ($1)
	RETURNING id, band_id, band_name, created_at 
	`
	var newBand models.Band

	err := DB.QueryRow(
		query,
		bandName,
	).Scan(
		&newBand.ID,
		&newBand.BandID,
		&newBand.Name,
		&newBand.CreatedAt,
	)
	if err != nil {
		fmt.Println("Unable to create band")
		return models.Band{}, err
	}
	return newBand, nil
}

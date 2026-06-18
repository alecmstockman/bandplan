package database

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateBandMembersTable(db *sql.DB) error {
	log.Println("- CreateBandMembersTable")
	query := `
	CREATE TABLE IF NOT EXISTS band_members (
		id SERIAL PRIMARY KEY, 
		band_id TEXT NOT NULL REFERENCES bands(band_id), 
		user_id TEXT NOT NULL REFERENCES users(user_id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Unable to create or load band members table")
		log.Fatal(err)
	}
	return nil
}

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

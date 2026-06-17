package database

import (
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
		band_email TEXT,
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

package database

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateBandsTable(db *sql.DB) error {
	fmt.Println("CreateBandsTable")
	query := `
	CREATE TABLE IF NOT EXITS bands (
		id SERIAL PRIMARY KEY,
		band_id TEXT NOT NULL,
		band_name TEXT NOT NULL, 
		band_email TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		updated_at NULL, 
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

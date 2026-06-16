package database

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateBandMembersTabe(db *sql.DB) error {
	fmt.Println("CreateBandMembersTable")

	query := `
	CREATE TABLE IF NOT EXISTS band_members (
		id SERIAL PRIMARY KEY, 
		band_id TEXT NOT NULL, 
		user_id TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'Member',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		
		FOREIGN KEY (band_id) REFERENCES bands.band_id,
		FOREIGN KEY (user_id) REFERENCES users.user_id
	)
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

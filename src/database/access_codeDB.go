package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/google/uuid"
)

func CreateAccessCodesTable(db *sql.DB) error {
	log.Println("- CreateAccessCodesTable")
	query := `
	CREATE TABLE IF NOT EXISTS access_codes (
		id SERIAL PRIMARY KEY,
		invite_id TEXT NOT NULL,
		code_hash TEXT NOT NULL,
		band_id TEXT NOT NULL REFERENCES bands(band_id),
		created_by TEXT NOT NULL REFERENCES users(user_id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		used_at TIMESTAMP,
		expires_at TIMESTAMP
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("   Unable to create Access Code Table: ", err)
		log.Fatal(err)
	}
	return nil
}

func AccessCodesTablesCreateCode(bandID string, userID string) (string, error) {
	log.Println("- CreateAccessCode")

	inviteID := uuid.NewString()

	code := uuid.NewString()
	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(1 * time.Hour)

	query := `
	INSERT INTO access_code
		code_hash, 
		band_id, 
		created_by,
		expires_at
	VALUES ($1, $2, $3)
	`
	_, err := DB.Exec(
		query,
		inviteID,
		codeHash,
		bandID,
		userID,
		expiresAt,
	)

	if err != nil {
		log.Println("   Unable to create access token: ", err)
		return "", err
	}

	return code, nil
}

func AccessCodesTableValidateCode(code string) (string, error) {
	log.Println("- AccessCodesTableValidateCode")

	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	query := `
	SELECT 
		expires_at, 
		band_id
	FROM access_codes
	WHERE code_hash = $1
	`
	var expiresAt time.Time
	var bandID string

	err := DB.QueryRow(
		query,
		codeHash,
	).Scan(
		&expiresAt,
		&bandID,
	)

	if err != nil {
		log.Println("   Unable to validate access code: ", err)
		return "", err
	}

	if expiresAt.Before(time.Now()) {
		log.Println("   Access code is expired: ", err)
		return "", nil
	}
	return bandID, nil
}

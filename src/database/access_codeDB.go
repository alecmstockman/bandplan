package database

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
)

func AccessCodesTablesCreateCode(bandID string, userID string) (string, error) {
	inviteID := uuid.NewString()

	code := strings.ToUpper(uuid.NewString()[:13])
	code = code[0:4] + "-" + code[4:]
	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(1 * time.Hour).UTC()

	query := `
	INSERT INTO access_codes(
		invite_id,
		code_hash, 
		band_id, 
		created_by,
		expires_at
	)
	VALUES ($1, $2, $3, $4, $5)
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
		return "", err
	}

	return code, nil
}

func AccessCodesTableValidateCode(code string) (string, error) {
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
		return "", err
	}

	if expiresAt.Before(time.Now().UTC()) {
		return "", nil
	}
	return bandID, nil
}

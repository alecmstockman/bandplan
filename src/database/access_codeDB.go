package database

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func AccessCodesTablesCreateCode(bandID string, userID string) (string, error) {
	fmt.Println("\n - AccessCodesTablesCreateCode")
	inviteID := uuid.NewString()

	code := strings.ToUpper(uuid.NewString()[:13])
	code = code[0:4] + "-" + code[4:]
	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	fmt.Println("create access code code: ", code)
	fmt.Println("create access code hash: ", hash)
	fmt.Println("create access code hash: ", codeHash)

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

func AccessCodesTableValidateCodeReturnBandID(code string) (string, error) {

	if code == "" {
		return "", errors.New("no access code provided")
	}

	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	fmt.Println("\n\n +++++++++++++ Access Code Hash: ", codeHash)
	fmt.Println("\n\n")

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

package database

import (
	"context"
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

func AccessCodesTableValidateCodeReturnBandID(ctx context.Context, code string) (string, error) {

	if code == "" {
		return "", errors.New("no access code provided")
	}

	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	query := `
		SELECT band_id
		FROM access_codes
		WHERE code_hash = $1
			AND expires_at > NOW()
	`

	var bandID string

	err := DB.QueryRowContext(ctx, query, codeHash).Scan(&bandID)
	if err != nil {
		return "", err
	}

	return bandID, nil
}

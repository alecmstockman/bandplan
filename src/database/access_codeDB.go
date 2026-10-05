package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrAccessCodeExpired = errors.New("access code expired")

func AccessCodesTablesCreateCode(bandID, userID, codeHash string, expiresAt time.Time) error {

	inviteID := uuid.NewString()

	query := `
		INSERT INTO access_codes(
			invite_id,
			code_hash, 
			band_id, 
			created_by,
			expires_at
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
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
		return err
	}

	return nil
}

func AccessCodesTableValidateCodeReturnBandID(ctx context.Context, code string) (string, error) {

	if code == "" {
		return "", errors.New("no access code provided")
	}

	query := `
		SELECT band_id, expires_at
		FROM access_codes
		WHERE code_hash = $1
	`

	var bandID string
	var expiresAt time.Time

	err := DB.QueryRowContext(ctx, query, code).Scan(&bandID, &expiresAt)
	if err != nil {
		return "", err
	}
	if !expiresAt.After(time.Now().UTC()) {
		return "", ErrAccessCodeExpired
	}

	return bandID, nil
}

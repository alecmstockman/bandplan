package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func BandEmailVerificationTableCreateCode(ctx context.Context, tx *sql.Tx, bandID, email, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO band_email_verification_tokens (
			band_id,
			email,
			token_hash,
			expires_at
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (band_id) DO UPDATE SET
			email = EXCLUDED.email,
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_at = CURRENT_TIMESTAMP
		WHERE band_email_verification_tokens.created_at <= CURRENT_TIMESTAMP - INTERVAL '60 seconds'
	`

	result, err := tx.ExecContext(ctx, query, bandID, email, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("BandEmailVerificationTableCreateCode - create code: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("BandEmailVerificationTableCreateCode - count affected verification tokens: %w", err)
	}
	if rowsAffected == 0 {
		return ErrEmailVerificationRateLimited
	}

	return nil
}

func BandEmailVerificationTableValidateCode(ctx context.Context, tokenHash string) (string, error) {
	if tokenHash == "" {
		return "", errors.New("no band email verification code provided")
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("BandEmailVerificationTableValidateCode - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var bandID string
	var email string
	err = tx.QueryRowContext(ctx, `
		DELETE FROM band_email_verification_tokens
		WHERE token_hash = $1
			AND expires_at > CURRENT_TIMESTAMP
		RETURNING band_id, email
	`, tokenHash).Scan(&bandID, &email)
	if err != nil {
		return "", fmt.Errorf("BandEmailVerificationTableValidateCode - consume verification code: %w", err)
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE bands
		SET
			is_email_verified = TRUE,
			email_verified_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE band_id = $1
			AND email_address = $2
	`, bandID, email)
	if err != nil {
		return "", fmt.Errorf("BandEmailVerificationTableValidateCode - verify email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("BandEmailVerificationTableValidateCode - count updated bands: %w", err)
	}
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("BandEmailVerificationTableValidateCode - commit transaction: %w", err)
	}

	return bandID, nil
}

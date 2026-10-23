package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrEmailVerificationCodeExpired = errors.New("email verification code expired")
var ErrEmailVerificationRateLimited = errors.New("email verification rate limited")

func EmailVerificationTableCreateCode(ctx context.Context, tx *sql.Tx, userID, email, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO email_verification_tokens (
			user_id,
			email,
			token_hash,
			expires_at
		) VALUES (
			$1, $2, $3, $4
		)
		ON CONFLICT (user_id) DO UPDATE SET
			email = EXCLUDED.email,
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_at = CURRENT_TIMESTAMP
		WHERE email_verification_tokens.created_at <= CURRENT_TIMESTAMP - INTERVAL '60 seconds'
	`

	result, err := tx.ExecContext(ctx, query, userID, email, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("EmailVerificationTableCreateCode - create code:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("EmailVerificationTableCreateCode - count affected verification tokens:  %w", err)
	}
	if rowsAffected == 0 {
		return ErrEmailVerificationRateLimited
	}

	return nil
}

func EmailVerificationTableValidateCodeReturnUserID(ctx context.Context, tokenHash string) (string, error) {
	if tokenHash == "" {
		return "", errors.New("no email verification code provided")
	}

	tx, err := DB.Begin()
	if err != nil {
		return "", fmt.Errorf("EmailVerificationTableValidateCodeReturnUserID - begin transaction:  %w", err)
	}

	defer tx.Rollback()

	deleteQuery := `
		DELETE FROM email_verification_tokens
		WHERE token_hash = $1
			AND expires_at > NOW()
		RETURNING 
			user_id, 
			expires_at
	`

	var userID string
	var expiresAt time.Time

	err = tx.QueryRowContext(ctx, deleteQuery, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		return "", fmt.Errorf("EmailVerificationTableValidateCodeReturnUserID - consume verification code:  %w", err)
	}
	if !expiresAt.After(time.Now().UTC()) {
		return "", ErrEmailVerificationCodeExpired
	}

	query := `
		UPDATE users
		SET
			is_email_verified = TRUE,
			email_verified_at = COALESCE(email_verified_at, CURRENT_TIMESTAMP),
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`

	result, err := tx.Exec(query, userID)
	if err != nil {
		return "", fmt.Errorf("EmailVerificationTableValidateCodeReturnUserID - verify email:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("EmailVerificationTableValidateCodeReturnUserID - count updated users:  %w", err)
	}
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("EmailVerificationTableValidateCodeReturnUserID - commit transaction:  %w", err)
	}

	return userID, nil
}

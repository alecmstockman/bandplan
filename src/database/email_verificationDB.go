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

func PasswordResetTokensTableCreateToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {

	query := `
		INSERT INTO password_reset_tokens (
			user_id, 
			token_hash,
			expires_at
		) VALUES (
			$1, $2, $3
		) ON CONFLICT (user_id) DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			expires_at = EXCLUDED.expires_at,
			created_at = CURRENT_TIMESTAMP
		WHERE password_reset_tokens.created_at <= CURRENT_TIMESTAMP - INTERVAL '60 seconds'
	`

	result, err := DB.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("PasswordResetTokensTableCreateToken - create code:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("PasswordResetTokensTableCreateToken - count affected verification tokens:  %w", err)
	}
	if rowsAffected == 0 {
		return ErrEmailVerificationRateLimited
	}

	return nil
}

func PasswordResetTokenTableValidateCodeReturnUserID(ctx context.Context, tokenHash, passwordHash string) (string, error) {

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("PasswordResetTokenTableValidateCodeReturnUserID - begin transaction:  %w", err)
	}

	defer tx.Rollback()

	deleteQuery := `
		DELETE FROM password_reset_tokens
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
		return "", fmt.Errorf("PasswordResetTokenTableValidateCodeReturnUserID - consume reset code:  %w", err)
	}
	if !expiresAt.After(time.Now().UTC()) {
		return "", fmt.Errorf("PasswordResetTokenTableValidateCodeReturnUserID - reset code expired:  %w", err)
	}

	query := `
		UPDATE users
		SET
			password_hash = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`

	result, err := tx.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return "", fmt.Errorf("PasswordResetTokenTableValidateCodeReturnUserID - verify email:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("PasswordResetTokenTableValidateCodeReturnUserID - count updated users:  %w", err)
	}
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("PasswordResetTokenTableValidateCodeReturnUserID - commit transaction:  %w", err)
	}

	return userID, nil
}

package database

import (
	"context"
	"errors"
	"time"
)

var ErrEmailVerificationCodeExpired = errors.New("email verification code expired")

func EmailVerificationTableCreateCode(userID, email, tokenHash string, expiresAt time.Time) error {
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
	`

	_, err := DB.Exec(query, userID, email, tokenHash, expiresAt)
	return err
}

func EmailVerificationTableValidateCodeReturnUserID(ctx context.Context, tokenHash string) (string, error) {
	if tokenHash == "" {
		return "", errors.New("no email verification code provided")
	}

	// query := `
	// 	SELECT user_id, expires_at
	// 	FROM email_verification_tokens
	// 	WHERE token_hash = $1
	// `

	query := `
		DELETE FROM email_verification_tokens
		WHERE token_hash = $1
			AND expires_at > NOW()
		RETURNING 
			user_id, 
			expires_at
	`

	var userID string
	var expiresAt time.Time

	err := DB.QueryRowContext(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		return "", err
	}
	if !expiresAt.After(time.Now().UTC()) {
		return "", ErrEmailVerificationCodeExpired
	}

	return userID, nil
}

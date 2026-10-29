-- +goose Up
CREATE TABLE band_email_verification_tokens (
    band_id TEXT PRIMARY KEY REFERENCES bands(band_id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE band_email_verification_tokens;

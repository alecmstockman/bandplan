-- +goose Up
CREATE TABLE user_registrations (
    id SERIAL PRIMARY KEY,
    user_registration_id TEXT NOT NULL UNIQUE,
    access_code_hash TEXT UNIQUE,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    display_name TEXT NOT NULL,

    band_id TEXT REFERENCES bands(band_id) ON DELETE SET NULL,
    timezone TEXT NOT NULL,

    email TEXT,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,

    password_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '24 hours',

    CHECK (NOT email_verified OR email IS NOT NULL),
    CHECK (password_hash IS NULL OR password_hash <> '')
);

CREATE UNIQUE INDEX idx_user_registrations_email
ON user_registrations (LOWER(email))
WHERE email IS NOT NULL;

-- +goose Down
DROP TABLE user_registration;

-- +goose Up
CREATE TABLE user_registration (
    id SERIAL PRIMARY KEY,
    user_registration_id TEXT NOT NULL UNIQUE,
    access_code_hash TEXT UNIQUE,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    slug TEXT UNIQUE,

    timezone TEXT,

    email TEXT NOT NULL UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,

    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
)

-- +goose Down
DROP TABLE user_registration;



-- +goose Up

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    slug TEXT UNIQUE,

    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,

    profile_image_id TEXT,
    profile_image_path TEXT,

    timezone TEXT,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,

    last_login TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

DROP TABLE users;
-- +goose Up
ALTER TABLE bands
    ADD COLUMN profile_image_id TEXT,
    ADD COLUMN profile_image_path TEXT,

    ADD COLUMN email_address TEXT UNIQUE,
    ADD COLUMN is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN email_verified_at TIMESTAMPTZ,
    ADD COLUMN timezone TEXT,

    ADD COLUMN city TEXT,
    ADD COLUMN state TEXT,
    ADD COLUMN country TEXT;

-- +goose Down
ALTER TABLE bands
    DROP COLUMN profile_image_id,
    DROP COLUMN profile_image_path,

    DROP COLUMN email_address,
    DROP COLUMN is_email_verified,
    DROP COLUMN email_verified_at,
    DROP COLUMN timezone,

    DROP COLUMN city,
    DROP COLUMN state,
    DROP COLUMN country;

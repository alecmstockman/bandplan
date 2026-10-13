-- +goose Up
ALTER TABLE band_members
    ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE band_members
    DROP COLUMN is_admin;


-- +goose Up
ALTER TABLE setlists
ADD COLUMN slug TEXT,
ADD COLUMN explicit BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN image_id TEXT;

-- +goose Down

ALTER TABLE setlists
DROP COLUMN image_id,
DROP COLUMN explicit, 
DROP COLUMN slug,

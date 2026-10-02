-- +goose Up
ALTER TABLE users
ADD COLUMN legal_accepted BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN legal_accepted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN legal_accepted,
DROP COLUMN legal_accepted_at;

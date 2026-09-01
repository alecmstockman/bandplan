-- +goose Up
ALTER TABLE message_reactions
ALTER COLUMN created_at TYPE TIMESTAMPTZ;

-- +goose Down
ALTER TABLE message_reactions
ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE;

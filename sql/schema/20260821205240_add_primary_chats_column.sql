-- +goose Up
ALTER TABLE chats
ADD COLUMN is_primary BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE chats
DROP COLUMN is_primary;

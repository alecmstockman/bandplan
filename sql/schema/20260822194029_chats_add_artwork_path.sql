-- +goose Up
ALTER TABLE chats
ADD COLUMN image_id TEXT,
ADD COLUMN image_path TEXT;

-- +goose Down
ALTER TABLE chats
DROP COLUMN IF EXISTS image_id,
DROP COLUMN IF EXISTS image_path;

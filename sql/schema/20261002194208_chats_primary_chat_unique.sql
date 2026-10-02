-- +goose Up
CREATE UNIQUE INDEX chats_one_primary_per_band
ON chats (band_id)
WHERE is_primary = true;

-- +goose Down
DROP INDEX chats_one_primary_per_band;
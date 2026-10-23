-- +goose Up

ALTER TABLE songs
ADD COLUMN itunes_track_id BIGINT;

-- +goose Down

ALTER TABLE songs
DROP COLUMN IF EXISTS itunes_track_id;

-- +goose Up
ALTER TABLE events
ADD COLUMN sound_check_time TIMESTAMPTZ;

-- +goose Down
ALTER TABLE events
DROP COLUMN sound_check_time;

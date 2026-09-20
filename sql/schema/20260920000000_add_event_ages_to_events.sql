-- +goose Up
ALTER TABLE events
ADD COLUMN event_ages TEXT NOT NULL DEFAULT 'N/A'
CHECK (event_ages IN ('N/A', 'allages', '18plus', '21plus'));

-- +goose Down
ALTER TABLE events
DROP COLUMN event_ages;

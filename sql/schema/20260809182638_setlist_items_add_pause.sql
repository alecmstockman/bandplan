-- +goose Up
ALTER TABLE setlist_items
ADD COLUMN pause_after_seconds INT NOT NULL DEFAULT 0
CHECK (pause_after_seconds >= 0);

-- +goose Down
ALTER TABLE setlist_items
DROP COLUMN pause_after_seconds;

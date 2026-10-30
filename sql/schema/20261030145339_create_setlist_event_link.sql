-- +goose Up
ALTER TABLE setlists
ADD COLUMN event_id TEXT REFERENCES events(event_id) ON DELETE SET NULL;

CREATE INDEX idx_setlists_event_id ON setlists(event_id);

-- +goose Down
DROP INDEX idx_setlists_event_id;

ALTER TABLE setlists
DROP COLUMN event_id;

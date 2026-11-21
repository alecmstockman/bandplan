-- +goose Up
ALTER TABLE events
    ADD COLUMN event_status TEXT NOT NULL CHECK (
        event_status IN (
            'active',
            'completed',
            'cancelled'
        )
    ) DEFAULT 'active';

-- +goose Down
ALTER TABLE events
DROP COLUMN event_status;

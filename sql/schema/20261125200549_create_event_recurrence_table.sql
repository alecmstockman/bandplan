-- +goose Up
CREATE TABLE event_recurrence (
    id SERIAL PRIMARY KEY,
    band_id TEXT NOT NULL REFERENCES bands(band_id) ON DELETE CASCADE,

    recurrence_id TEXT NOT NULL UNIQUE,
    series_id TEXT NOT NULL,
    event_id TEXT NOT NULL UNIQUE
        REFERENCES events(event_id)
        ON DELETE CASCADE,

    recurrence TEXT NOT NULL CHECK (
        recurrence IN (
            'none',
            'daily',
            'weekly',
            'biweekly',
            'monthly',
            'yearly'
        )
    ),
    recurrence_count INT CHECK (recurrence_count > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id)
);

CREATE INDEX idx_event_recurrence_series_id
	ON event_recurrence(series_id);

CREATE INDEX idx_event_recurrence_band_id
	ON event_recurrence(band_id);

-- +goose Down
DROP TABLE event_recurrence;

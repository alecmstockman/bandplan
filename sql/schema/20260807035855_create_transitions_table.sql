

-- +goose Up

CREATE TABLE transitions (
    id SERIAL PRIMARY KEY,
    transition_id TEXT NOT NULL UNIQUE,
    band_id TEXT NOT NULL REFERENCES bands(band_id) ON DELETE CASCADE,

    title TEXT NOT NULL,
    title_slug TEXT NOT NULL,

    length_seconds int NOT NULL DEFAULT 0,
    bpm int NOT NULL DEFAULT 0,
    time_signature TEXT NOT NULL DEFAULT '',
    musical_key TEXT NOT NULL DEFAULT '',

    tuning TEXT NOT NULL DEFAULT '',
    capo TEXT NOT NULL DEFAULT '',

    explicit BOOLEAN NOT NULL DEFAULT FALSE,
    chords TEXT NOT NULL DEFAULT '',
    chart_link TEXT NOT NULL DEFAULT '',

    lyrics TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL,

    UNIQUE (band_id, title_slug)
);

CREATE INDEX idx_transitions_band_id
ON transitions (band_id);


-- +goose Down

DROP TABLE transitions;

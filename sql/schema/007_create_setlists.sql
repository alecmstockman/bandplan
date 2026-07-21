

-- +goose Up

CREATE TABLE setlists (
    id SERIAL PRIMARY KEY,
    setlist_id TEXT NOT NULL UNIQUE, 
    band_id TEXT NOT NULL REFERENCES bands(band_id),

    name TEXT NOT NULL,
    notes TEXT,
    artwork_path TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    UNIQUE(band_id, name)
);


-- +goose Down

DORP TABLE setlists;
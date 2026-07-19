

-- +goose Up

CREATE TABLE band_members (
    id SERIAL PRIMARY KEY,

    band_id TEXT NOT NULL REFERENCES bands(band_id),
    user_id TEXT NOT NULL REFERENCES users(user_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    UNIQUE (band_id, user_id)
);

-- +goose Down

DROP TABLE band_members;
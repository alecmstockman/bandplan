

-- +goose Up

CREATE TABLE access_codes (
    id SERIAL PRIMARY KEY,

    invite_id TEXT NOT NULL UNIQUE,
    code_hash TEXT NOT NULL UNIQUE,

    band_id TEXT NOT NULL REFERENCES bands(band_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    
    used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE access_codes;
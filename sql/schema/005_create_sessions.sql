

-- +goose Up

CREATE TABLE sessions (
	id SERIAL PRIMARY KEY,
	
	user_id TEXT NOT NULL REFERENCES users(user_id),
	band_id TEXT REFERENCES bands(band_id),

	token TEXT NOT NULL UNIQUE,

	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMPTZ NOT NULL
);


-- +goose Down

DROP TABLE sessions;
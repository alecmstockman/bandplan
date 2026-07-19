

-- +goose Up

CREATE TABLE messages (
	id SERIAL PRIMARY KEY,
	message_id TEXT NOT NULL UNIQUE,

	band_id TEXT NOT NULL REFERENCES bands(band_id),
	user_id TEXT NOT NULL REFERENCES users(user_id),

	body TEXT NOT NULL,
	
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	edited_at TIMESTAMP
);

-- +goose Down

DROP TABLE messages;
-- +goose Up
CREATE TABLE notifications
	id SERIAL PRIMARY KEY,
	notification_id TEXT NOT NULL UNIQUE,

	band_id TEXT NOT NULL REFERENCES bands(band_id),
	user_id TEXT NOT NULL REFERENCES users(user_id),

	body TEXT NOT NULL,
	
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	edited_at TIMESTAMP

-- +goose Down
DROP TABLE notifications;

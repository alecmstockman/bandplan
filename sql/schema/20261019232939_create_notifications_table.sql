-- +goose Up
CREATE TABLE notifications (
	id SERIAL PRIMARY KEY,
	notification_id TEXT NOT NULL UNIQUE,

	band_id TEXT REFERENCES bands(band_id),
	user_id TEXT NOT NULL REFERENCES users(user_id),

	type TEXT NOT NULL,
	item_id TEXT NOT NULL, 
	body TEXT NOT NULL,
	
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by TEXT NOT NULL REFERENCES users(user_id),
	read_at TIMESTAMPTZ
);

CREATE INDEX notifications_user_band_created_idx
	ON notifications (user_id, band_id, created_at DESC, id DESC);

-- +goose Down
DROP TABLE notifications;

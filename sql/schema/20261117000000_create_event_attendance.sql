-- +goose Up
CREATE TABLE event_attendance (
    id SERIAL PRIMARY KEY,
    event_attendance_id TEXT NOT NULL UNIQUE,
    event_id TEXT NOT NULL REFERENCES events(event_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    status TEXT NOT NULL DEFAULT 'unseen',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (event_id, user_id),

    CONSTRAINT event_attendance_status_check
        CHECK (status IN ('unseen', 'maybe', 'confirmed', 'cannot'))
);

CREATE INDEX idx_event_attendance_event_id
ON event_attendance(event_id);

CREATE INDEX idx_event_attendance_user_id
ON event_attendance(user_id);

-- +goose Down
DROP TABLE event_attendance;

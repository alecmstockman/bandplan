-- +goose Up
ALTER TABLE messages
    ADD COLUMN is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN pinned_at TIMESTAMPTZ,
    ADD COLUMN pinned_by TEXT 
        REFERENCES users(user_id)
        ON DELETE SET NULL;

-- +goose Down
ALTER TABLE messages
    DROP COLUMN is_pinned,
    DROP COLUMN pinned_at,
    DROP COLUMN pinned_by;

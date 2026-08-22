-- +goose Up

CREATE TABLE chats (
    id SERIAL PRIMARY KEY, 
    chat_id TEXT NOT NULL UNIQUE,
    band_id TEXT NOT NULL REFERENCES bands(band_id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    slug TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL DEFAULT '',

    UNIQUE (band_id, slug)
);

CREATE TABLE chat_members (
    id SERIAL PRIMARY KEY,
    chat_id TEXT NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT,

    UNIQUE (chat_id, user_id)
);

CREATE INDEX idx_chats_band_id
ON chats(band_id);

CREATE INDEX idx_chat_members_user_id
ON chat_members(user_id);


-- +goose Down

DROP TABLE chat_members;
DROP TABLE chats;
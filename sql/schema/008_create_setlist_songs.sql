

-- +goose Up

CREATE TABLE setlist_songs (
    id SERIAL PRIMARY KEY,

    setlist_id TEXT NOT NULL 
        REFERENCES setlists(setlist_id)
        ON DELETE CASCADE,

    song_id TEXT NOT NULL 
        REFERENCES songs(song_id)
        ON DELETE CASCADE,

    position INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    UNIQUE (setlist_id, position)
);

-- +goose Down

DROP TABLE setlist_songs;
-- +goose Up

CREATE TABLE todo_lists (
    id SERIAL PRIMARY KEY,
    todo_list_id TEXT NOT NULL UNIQUE,
    user_id TEXT REFERENCES users(user_id),
    band_id TEXT REFERENCES bands(band_id),

    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    description TEXT,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    song_id TEXT REFERENCES songs(song_id) ON DELETE SET NULL,
    setlist_id TEXT REFERENCES setlists(setlist_id) ON DELETE SET NULL,
    event_id TEXT REFERENCES events(event_id) ON DELETE SET NULL,

    assigned_to TEXT REFERENCES users(user_id),
    due_date DATE,
    due_time TIME,
    due_timezone TEXT,

    is_complete BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    completed_by TEXT REFERENCES users(user_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    CONSTRAINT todo_lists_one_owner
        CHECK ((user_id IS NOT NULL) <> (band_id IS NOT NULL)),

    CHECK (
        (
            is_complete = FALSE
            AND completed_at IS NULL
            AND completed_by IS NULL
        )
        OR
        (
            is_complete = TRUE
            AND completed_at IS NOT NULL
            AND completed_by IS NOT NULL
        )
    )
);

CREATE UNIQUE INDEX todo_lists_one_primary_per_user
ON todo_lists(user_id)
WHERE is_primary = TRUE AND user_id IS NOT NULL;

CREATE UNIQUE INDEX todo_lists_one_primary_per_band
ON todo_lists(band_id)
WHERE is_primary = TRUE AND band_id IS NOT NULL;

CREATE TABLE todo_items (
    id SERIAL PRIMARY KEY,
    item_id TEXT NOT NULL UNIQUE,
    todo_list_id TEXT NOT NULL REFERENCES todo_lists(todo_list_id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    position INT NOT NULL,
    body TEXT,

    song_id TEXT REFERENCES songs(song_id) ON DELETE SET NULL,
    setlist_id TEXT REFERENCES setlists(setlist_id) ON DELETE SET NULL,
    event_id TEXT REFERENCES events(event_id) ON DELETE SET NULL,

    assigned_to TEXT REFERENCES users(user_id),
    due_date DATE,
    due_time TIME,
    due_timezone TEXT,

    is_complete BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    completed_by TEXT REFERENCES users(user_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    CONSTRAINT todo_position_unique
        UNIQUE (todo_list_id, position)
        DEFERRABLE INITIALLY IMMEDIATE,

    CHECK (position >= 0),

    CHECK (
        (
            is_complete = FALSE
            AND completed_at IS NULL
            AND completed_by IS NULL
        )
        OR
        (
            is_complete = TRUE
            AND completed_at IS NOT NULL
            AND completed_by IS NOT NULL
        )
    )
);

CREATE INDEX idx_todo_lists_band_id ON todo_lists(band_id);
CREATE INDEX idx_todo_lists_user_id ON todo_lists(user_id);


CREATE TABLE checklist_lists (
    id SERIAL PRIMARY KEY,
    checklist_id TEXT NOT NULL UNIQUE,
    user_id TEXT REFERENCES users(user_id),
    band_id TEXT REFERENCES bands(band_id),

    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    description TEXT,

    song_id TEXT REFERENCES songs(song_id) ON DELETE SET NULL,
    setlist_id TEXT REFERENCES setlists(setlist_id) ON DELETE SET NULL,
    event_id TEXT REFERENCES events(event_id) ON DELETE SET NULL,

    assigned_to TEXT REFERENCES users(user_id),
    due_date DATE,
    due_time TIME,
    due_timezone TEXT,

    is_complete BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    completed_by TEXT REFERENCES users(user_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    CONSTRAINT checklist_lists_one_owner
        CHECK ((user_id IS NOT NULL) <> (band_id IS NOT NULL)),

    CHECK (
        (
            is_complete = FALSE
            AND completed_at IS NULL
            AND completed_by IS NULL
        )
        OR
        (
            is_complete = TRUE
            AND completed_at IS NOT NULL
            AND completed_by IS NOT NULL
        )
    )
);

CREATE TABLE checklist_items (
    id SERIAL PRIMARY KEY,
    item_id TEXT NOT NULL UNIQUE,
    checklist_id TEXT NOT NULL REFERENCES checklist_lists(checklist_id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    position INT NOT NULL,
    body TEXT,

    song_id TEXT REFERENCES songs(song_id) ON DELETE SET NULL,
    setlist_id TEXT REFERENCES setlists(setlist_id) ON DELETE SET NULL,
    event_id TEXT REFERENCES events(event_id) ON DELETE SET NULL,

    assigned_to TEXT REFERENCES users(user_id),
    due_date DATE,
    due_time TIME,
    due_timezone TEXT,

    is_complete BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    completed_by TEXT REFERENCES users(user_id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    CONSTRAINT checklist_position_unique
        UNIQUE (checklist_id, position)
        DEFERRABLE INITIALLY IMMEDIATE,

    CHECK (position >= 0),

    CHECK (
        (
            is_complete = FALSE
            AND completed_at IS NULL
            AND completed_by IS NULL
        )
        OR
        (
            is_complete = TRUE
            AND completed_at IS NOT NULL
            AND completed_by IS NOT NULL
        )
    )
);

CREATE INDEX idx_checklist_lists_band_id ON checklist_lists(band_id);
CREATE INDEX idx_checklist_lists_user_id ON checklist_lists(user_id);


-- +goose Down

DROP TABLE checklist_items;
DROP TABLE checklist_lists;

DROP TABLE todo_items;
DROP TABLE todo_lists;


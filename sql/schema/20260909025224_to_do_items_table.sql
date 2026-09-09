-- +goose Up

CREATE TABLE to_do_lists (
    id SERIAL PRIMARY KEY,
    to_do_list_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,

    band_id TEXT REFERENCES bands(band_id)
        ON DELETE CASCADE,

    owner_user_id TEXT REFERENCES users(user_id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    CONSTRAINT to_do_lists_owner_check CHECK (
        (
            band_id IS NOT NULL
            AND owner_user_id IS NULL
        )
        OR
        (
            band_id IS NULL
            AND owner_user_id IS NOT NULL
        )
    )
);

CREATE TABLE to_do_items (
    id SERIAL PRIMARY KEY,
    item_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,

    to_do_list_id TEXT NOT NULL
        REFERENCES to_do_lists(to_do_list_id)
        ON DELETE CASCADE,

    is_complete BOOLEAN NOT NULL DEFAULT FALSE,
    body TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id)
);

CREATE INDEX idx_to_do_lists_band_id
ON to_do_lists(band_id);

CREATE INDEX idx_to_do_lists_owner_user_id
ON to_do_lists(owner_user_id);

CREATE INDEX idx_to_do_items_to_do_list_id
ON to_do_items(to_do_list_id);


-- +goose Down

DROP TABLE to_do_items;

DROP TABLE to_do_lists;
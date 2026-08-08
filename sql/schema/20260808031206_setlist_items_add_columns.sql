-- +goose Up

CREATE TABLE breaks (
    id SERIAL PRIMARY KEY,
    break_id TEXT NOT NULL UNIQUE,
    band_id TEXT NOT NULL REFERENCES bands(band_id),

    title TEXT NOT NULL,
    title_slug TEXT NOT NULL, 

    length_seconds INT NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '', 

    link_one TEXT NOT NULL DEFAULT '',
    link_two TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id)
);

ALTER TABLE setlist_items
ALTER COLUMN song_id DROP NOT NULL,
ADD COLUMN item_type TEXT NOT NULL DEFAULT 'song',
ADD COLUMN transition_id TEXT
    REFERENCES transitions(transition_id)
        ON DELETE CASCADE,
ADD COLUMN break_id TEXT
    REFERENCES breaks(break_id)
        ON DELETE CASCADE;

ALTER TABLE setlist_items
ADD CONSTRAINT setlist_items_exactly_one_item_check
CHECK (
    (song_id IS NOT NULL)::int +
    (transition_id IS NOT NULL)::int +
    (break_id IS NOT NULL)::int = 1
);

ALTER TABLE setlist_items
ADD CONSTRAINT setlist_items_item_type_check
CHECK (
    (
        item_type = 'song'
        AND song_id IS NOT NULL
        AND transition_id IS NULL
        AND break_id IS NULL
    )
    OR 
    (
        item_type = 'transition'
        AND transition_id IS NOT NULL
        AND song_id IS NULL
        AND break_id IS NULL
    )
    OR 
    (
        item_type = 'break'
        AND break_id IS NOT NULL
        AND song_id IS NULL
        AND transition_id IS NULL
    )
);

-- +goose Down

ALTER TABLE setlist_items
DROP CONSTRAINT IF EXISTS setlist_items_item_type_check;

ALTER TABLE setlist_items
DROP CONSTRAINT IF EXISTS setlist_items_exactly_one_item_check;

ALTER TABLE setlist_items
DROP COLUMN IF EXISTS break_id;

ALTER TABLE setlist_items
DROP COLUMN IF EXISTS transition_id;

ALTER TABLE setlist_items
DROP COLUMN IF EXISTS item_type;

ALTER TABLE setlist_items
ALTER COLUMN song_id SET NOT NULL;

DROP TABLE breaks;

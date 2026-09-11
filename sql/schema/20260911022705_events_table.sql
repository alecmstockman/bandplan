-- +goose Up
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_id TEXT NOT NULL UNIQUE,
    band_id TEXT NOT NULL 
        REFERENCES bands(band_id) 
        ON DELETE CASCADE,
    
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    image_id TEXT NOT NULL DEFAULT '',
    image_path TEXT NOT NULL DEFAULT '', 

    event_date DATE NOT NULL,
    event_type TEXT NOT NULL CHECK (
        event_type IN (
            'show', 'gig', 'festival', 'rehearsal', 'practice',
            'writing', 'recording', 'meeting', 'photos', 'press', 'other'
        )
    ),
    recurrence TEXT NOT NULL,

    location TEXT,
    address TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    time_zone TEXT NOT NULL,

    set_location TEXT,
    load_in_time TIMESTAMPTZ, 
    load_in_instructions TEXT,
    set_time TIMESTAMPTZ,
    set_length_seconds INTEGER 
        CHECK (set_length_seconds >= 0),

    venue_name TEXT,
    address_one TEXT,
    address_two TEXT,
    city TEXT,
    state TEXT,
    zip_code TEXT,

    presale_ticket_price NUMERIC(10,2)
        CHECK (presale_ticket_price >= 0),
    ticket_price NUMERIC(10,2)
        CHECK (ticket_price >= 0),
    ticket_link TEXT,

    setlist_id TEXT 
        REFERENCES setlists(setlist_id)
        ON DELETE SET NULL,
    notes TEXT,

    link_one_name TEXT,
    link_one TEXT,
    link_two_name TEXT,
    link_two TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id),

    UNIQUE (band_id, slug)
);

CREATE INDEX idx_events_setlist_id
ON events(setlist_id);

CREATE INDEX idx_events_band_date
ON events(band_id, event_date);

CREATE INDEX idx_events_event_date
ON events(event_date);

-- +goose Down
DROP TABLE events;



-- +goose Up

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    song_id TEXT NOT NULL UNIQUE,
    band_id TEXT NOT NULL REFERENCES bands(band_id),
    
    title TEXT NOT NULL,
    title_slug TEXT,
    album_title TEXT,
    album_id TEXT,
    album_slug TEXT,
    artist_name TEXT NOT NULL,
    artist_id TEXT,
    artist_slug TEXT,

    artwork_id TEXT,
    artwork_path TEXT,
    release_date DATE,
    genre TEXT,

    recording_bpm INTEGER,
    live_bpm INTEGER,
    time_signature TEXT,
    original_key TEXT,
    live_key TEXT,
    tuning TEXT,
    capo TEXT,
    length_seconds INTEGER,

    status TEXT,
    explicit BOOLEAN NOT NULL DEFAULT FALSE,
    is_cover BOOLEAN NOT NULL DEFAULT FALSE,
    chords TEXT, 
    chart_link TEXT,

    spotify_link TEXT,
    apple_music_link TEXT,
    youtube_link TEXT,
    amazon_music_link TEXT,
    pandora_link TEXT,
    deezer_link TEXT,
    tidal_link TEXT,
    other_link TEXT,

    lyrics TEXT,
    description TEXT,
    notes TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL REFERENCES users(user_id),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT REFERENCES users(user_id)
);

-- +goose Down

DROP TABLE songs;
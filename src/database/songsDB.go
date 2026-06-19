package database

import (
	"database/sql"
	"log"
)

func CreateSongsTable(db *sql.DB) error {
	log.Println("- CreateSongsTable")
	query := `
	CREATE TABLE IF NOT EXITS songs (
		id SERIAL PRIMARY KEY,
		song_id TEXT,
		name TEXT,
		album_name TEXT,
		band_id TEXT,
		genre TEXT,

		musical_key TEXT,
		tuning TEXT,
		recording_bpm INTEGER,
		live_bpm INTEGER,
		length_seconds INTEGER,

		release_date DATE,
		lyrics TEXT,
		spotify_link TEXT,
		apple_music_link TEXT,
		youtube_link TEXT,
		other_link TEXT,

		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Println("   Unable to create or load songs table: ", err)
		log.Fatal(err)
	}
	return nil
}

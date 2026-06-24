package database

import (
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func CreateSongsTable(db *sql.DB) error {
	log.Println("- CreateSongsTable")
	query := `
	CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		song_id TEXT,
		title TEXT,
		album_title TEXT,
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

func SongsTableCreateSong(song models.Song) (models.Song, error) {
	fmt.Println("-------------------------------------------------------------")
	log.Println("- SongsTableCreateSong")
	fmt.Println("-------------------------------------------------------------")

	songID := uuid.New().String()

	var newSong models.Song

	query := `
	INSERT INTO songs(
		song_id,
		title,
		album_title,
		band_id,
		genre,

		musical_key,
		tuning,
		recording_bpm,
		live_bpm,
		length_seconds,

		release_date,
		spotify_link,
		apple_music_link,
		youtube_link,
		other_link,

		lyrics,
		notes

	)
	VALUES (
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16,
		$17
	)
	RETURNING 
		id, 

		song_id,
		title,
		album_title,
		band_id,
		genre,

		musical_key,
		tuning,
		recording_bpm,
		live_bpm,
		length_seconds,

		release_date,
		spotify_link,
		apple_music_link,
		youtube_link,
		other_link,
		
		lyrics,
		notes,
		created_at,
		updated_at
	`

	err := DB.QueryRow(
		query,
		songID,
		song.Title,
		song.AlbumTitle,
		song.BandID,
		song.Genre,

		song.MusicalKey,
		song.Tuning,
		song.RecordingBPM,
		song.LiveBPM,
		song.LiveBPM,

		song.ReleaseDate,
		song.SpotifyLink,
		song.AppleMusicLink,
		song.YouTubeLink,
		song.OtherLink,

		song.Lyrics,
		song.Notes,
	).Scan(
		&newSong.ID,
		&newSong.SongID,
		&newSong.Title,
		&newSong.AlbumTitle,
		&newSong.BandID,
		&newSong.Genre,

		&newSong.MusicalKey,
		&newSong.Tuning,
		&newSong.RecordingBPM,
		&newSong.LiveBPM,
		&newSong.LengthSeconds,

		&newSong.ReleaseDate,
		&newSong.SpotifyLink,
		&newSong.AppleMusicLink,
		&newSong.YouTubeLink,
		&newSong.OtherLink,

		&newSong.Lyrics,
		&newSong.Notes,
		&newSong.CreatedAt,
		&newSong.UpdatedAt,
	)

	if err != nil {
		log.Println("   Unable to save song to db: ", err)
		return models.Song{}, err
	}

	return newSong, nil
}

func SongsTableGetAllSongsByBandID(bandID string) ([]models.Song, error) {
	log.Println("- SongsTableGetAllSongsByBandID")

	query := `
	SELECT
		id,
		song_id,
		title,
		album_title,
		band_id,
		genre,

		musical_key,
		tuning,
		recording_bpm,
		live_bpm,
		length_seconds,

		release_date,
		spotify_link,
		apple_music_link,
		youtube_link,
		other_link,

		lyrics,
		notes,
		created_at,
		updated_at
	FROM songs
	WHERE band_id = $1
	ORDER BY title ASC
	`

	rows, err := DB.Query(query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song

	for rows.Next() {
		var song models.Song

		err := rows.Scan(
			&song.ID,
			&song.SongID,
			&song.Title,
			&song.AlbumTitle,
			&song.BandID,
			&song.Genre,

			&song.MusicalKey,
			&song.Tuning,
			&song.RecordingBPM,
			&song.LiveBPM,
			&song.LengthSeconds,

			&song.ReleaseDate,
			&song.SpotifyLink,
			&song.AppleMusicLink,
			&song.YouTubeLink,
			&song.OtherLink,

			&song.Lyrics,
			&song.Notes,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func SongsTableSearchByBandID(bandID string, query string) ([]models.Song, error) {
	rows, err := DB.Query(`
		SELECT
			id, song_id, title, album_title, band_id, genre,
			musical_key, tuning, recording_bpm, live_bpm, length_seconds,
			release_date, lyrics, spotify_link, apple_music_link, youtube_link, other_link,
			notes, created_at, updated_at
		FROM songs
		WHERE band_id = $1
		AND (
			$2 = ''
			OR title ILIKE '%' || $2 || '%'
			OR album_title ILIKE '%' || $2 || '%'
			OR genre ILIKE '%' || $2 || '%'
			OR musical_key ILIKE '%' || $2 || '%'
			OR tuning ILIKE '%' || $2 || '%'
			OR recording_bpm::text ILIKE '%' || $2 || '%'
			OR live_bpm::text ILIKE '%' || $2 || '%'
			OR length_seconds::text ILIKE '%' || $2 || '%'
			OR lyrics ILIKE '%' || $2 || '%'
			OR notes ILIKE '%' || $2 || '%'
		)
		ORDER BY title ASC
	`, bandID, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song

	for rows.Next() {
		var song models.Song

		err := rows.Scan(
			&song.ID,
			&song.SongID,
			&song.Title,
			&song.AlbumTitle,
			&song.BandID,
			&song.Genre,
			&song.MusicalKey,
			&song.Tuning,
			&song.RecordingBPM,
			&song.LiveBPM,
			&song.LengthSeconds,
			&song.ReleaseDate,
			&song.Lyrics,
			&song.SpotifyLink,
			&song.AppleMusicLink,
			&song.YouTubeLink,
			&song.OtherLink,
			&song.Notes,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, rows.Err()
}

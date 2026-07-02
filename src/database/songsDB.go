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
		song_id TEXT UNIQUE NOT NULL,
		band_id TEXT,
		
		title TEXT,
		title_slug TEXT,
		album_title TEXT,
		album_id TEXT,
		album_slug TEXT,
		artist_name TEXT,
		artist_id TEXT,
		artist_slug TEXT,

		artwork_id TEXT,
		artwork_path TEXT,
		release_date DATE,
		genre TEXT,

		recording_bpm INTEGER,
		live_bpm INTEGER,
		musical_key TEXT,
		tuning TEXT,
		capo TEXT,
		length_seconds INTEGER,

		status TEXT,
		explicit BOOLEAN NOT NULL DEFAULT FALSE,
		is_cover BOOLEAN NOT NULL DEFAULT FALSE,

		spotify_link TEXT,
		apple_music_link TEXT,
		youtube_link TEXT,
		amazon_music_link TEXT,
		pandora_link TEXT,
		deezer_link TEXT,
		other_link TEXT,

		lyrics TEXT,
		description TEXT,
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
	fmt.Println("\n\n-------------------------------------------------------------")
	log.Println("- SongsTableCreateSong")

	fmt.Println(song.LengthSeconds)

	songID := uuid.New().String()

	var newSong models.Song

	query := `
	INSERT INTO songs(
		song_id,
		band_id,

		title,
		title_slug,
		album_title,
		album_id,
		album_slug,
		artist_name,
		artist_id,
		artist_slug,

		artwork_id,
		artwork_path,
		release_date,
		genre,

		recording_bpm,
		live_bpm,
		musical_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		other_link,

		lyrics,
		description,
		notes
	)
	VALUES (
		$1, $2,
		$3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20,
		$21, $22, $23,
		$24, $25, $26, $27, $28, $29, $30,
		$31, $32, $33
	)
	RETURNING 
		id, 
		song_id,
		band_id,

		title,
		title_slug,
		album_title,
		album_id, 
		album_slug,
		artist_name,
		artist_id,
		artist_slug,

		artwork_id,
		artwork_path,
		release_date,
		genre,

		recording_bpm,
		live_bpm,
		musical_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		other_link,

		lyrics,
		description,
		notes,
		created_at,
		updated_at
	`

	err := DB.QueryRow(
		query,
		songID,
		song.BandID,

		song.Title,
		song.TitleSlug,
		song.AlbumTitle,
		song.AlbumID,
		song.AlbumSlug,
		song.ArtistName,
		song.ArtistID,
		song.ArtistSlug,

		song.ArtworkID,
		song.ArtworkPath,
		song.ReleaseDate,
		song.Genre,

		song.RecordingBPM,
		song.LiveBPM,
		song.MusicalKey,
		song.Tuning,
		song.Capo,
		song.LengthSeconds,

		song.Status,
		song.Explicit,
		song.IsCover,

		song.SpotifyLink,
		song.AppleMusicLink,
		song.YouTubeLink,
		song.AmazonMusicLink,
		song.PandoraLink,
		song.DeezerLink,
		song.OtherLink,

		song.Lyrics,
		song.Description,
		song.Notes,
	).Scan(
		&newSong.ID,
		&newSong.SongID,
		&newSong.BandID,

		&newSong.Title,
		&newSong.TitleSlug,
		&newSong.AlbumTitle,
		&newSong.AlbumID,
		&newSong.AlbumSlug,
		&newSong.ArtistName,
		&newSong.ArtistID,
		&newSong.ArtistSlug,

		&newSong.ArtworkID,
		&newSong.ArtworkPath,
		&newSong.ReleaseDate,
		&newSong.Genre,

		&newSong.RecordingBPM,
		&newSong.LiveBPM,
		&newSong.MusicalKey,
		&newSong.Tuning,
		&newSong.Capo,
		&newSong.LengthSeconds,

		&newSong.Status,
		&newSong.Explicit,
		&newSong.IsCover,

		&newSong.SpotifyLink,
		&newSong.AppleMusicLink,
		&newSong.YouTubeLink,
		&newSong.AmazonMusicLink,
		&newSong.PandoraLink,
		&newSong.DeezerLink,
		&newSong.OtherLink,

		&newSong.Lyrics,
		&newSong.Description,
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
		band_id,

		title,
		title_slug,
		album_title,
		album_id, 
		album_slug,
		artist_name,
		artist_id,
		artist_slug,

		artwork_id,
		artwork_path,
		release_date,
		genre,

		recording_bpm,
		live_bpm,
		musical_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		other_link,

		lyrics,
		description,
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
			&song.BandID,

			&song.Title,
			&song.TitleSlug,
			&song.AlbumTitle,
			&song.AlbumID,
			&song.AlbumSlug,
			&song.ArtistID,
			&song.ArtistSlug,

			&song.ArtistID,
			&song.ArtworkPath,
			&song.ReleaseDate,
			&song.Genre,

			&song.RecordingBPM,
			&song.LiveBPM,
			&song.MusicalKey,
			&song.Tuning,
			&song.Capo,
			&song.LengthSeconds,

			&song.Status,
			&song.Explicit,
			&song.IsCover,

			&song.ReleaseDate,
			&song.SpotifyLink,
			&song.AppleMusicLink,
			&song.YouTubeLink,
			&song.AmazonMusicLink,
			&song.PandoraLink,
			&song.DeezerLink,
			&song.OtherLink,

			&song.Lyrics,
			&song.Description,
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
			id,
			song_id,
			band_id,

			title,
			title_slug,
			album_title,
			album_id,
			album_slug,
			artist_name,
			artist_id,
			artist_slug,

			artwork_id,
			artwork_path,
			release_date,
			genre,

			recording_bpm,
			live_bpm,
			musical_key,
			tuning,
			capo,
			length_seconds,

			status,
			explicit,
			is_cover,

			spotify_link,
			apple_music_link,
			youtube_link,
			amazon_music_link,
			pandora_link,
			deezer_link,
			other_link,

			lyrics,
			description,
			notes,
			created_at,
			updated_at
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
			&song.BandID,

			&song.Title,
			&song.TitleSlug,
			&song.AlbumTitle,
			&song.AlbumID,
			&song.AlbumSlug,
			&song.ArtistID,
			&song.ArtistSlug,

			&song.ArtistID,
			&song.ArtworkPath,
			&song.ReleaseDate,
			&song.Genre,

			&song.RecordingBPM,
			&song.LiveBPM,
			&song.MusicalKey,
			&song.Tuning,
			&song.Capo,
			&song.LengthSeconds,

			&song.Status,
			&song.Explicit,
			&song.IsCover,

			&song.ReleaseDate,
			&song.SpotifyLink,
			&song.AppleMusicLink,
			&song.YouTubeLink,
			&song.AmazonMusicLink,
			&song.PandoraLink,
			&song.DeezerLink,
			&song.OtherLink,

			&song.Lyrics,
			&song.Description,
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

func SongsTableGetSongBySongID(songID string) (models.Song, error) {
	log.Println("- SongsTableGetSongBySongID")

	query := `
	SELECT * 
	FROM songs
	WHERE song_id = $1
	`

	var song models.Song

	err := DB.QueryRow(query, songID).Scan(
		&song.ID,
		&song.SongID,
		&song.BandID,

		&song.Title,
		&song.TitleSlug,
		&song.AlbumTitle,
		&song.AlbumID,
		&song.AlbumSlug,
		&song.ArtistID,
		&song.ArtistSlug,

		&song.ArtistID,
		&song.ArtworkPath,
		&song.ReleaseDate,
		&song.Genre,

		&song.RecordingBPM,
		&song.LiveBPM,
		&song.MusicalKey,
		&song.Tuning,
		&song.Capo,
		&song.LengthSeconds,

		&song.Status,
		&song.Explicit,
		&song.IsCover,

		&song.ReleaseDate,
		&song.SpotifyLink,
		&song.AppleMusicLink,
		&song.YouTubeLink,
		&song.AmazonMusicLink,
		&song.PandoraLink,
		&song.DeezerLink,
		&song.OtherLink,

		&song.Lyrics,
		&song.Description,
		&song.Notes,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
	if err != nil {
		log.Println("   Unable to get song from songs db: ", err)
		return models.Song{}, nil
	}

	return song, nil
}

func SongsTableUpdateSong(song models.Song) error {
	log.Println("- SongsTableUpdateSong")

	query := `
		UPDATE songs
		SET
			title = $1,
			title_slug = $2,

			album_title = $3,
			album_id = $4,
			album_slug = $5,

			artist_name = $6,
			artist_id = $7,
			artist_slug = $8,

			artwork_id = $9,
			artwork_path = $10,

			release_date = $11,
			genre = $12,

			recording_bpm = $13,
			live_bpm = $14,
			musical_key = $15,
			tuning = $16,
			capo = $17,
			length_seconds = $18,

			status = $19,
			explicit = $20,
			is_cover = $21,

			spotify_link = $22,
			apple_music_link = $23,
			youtube_link = $24,
			amazon_music_link = $25,
			pandora_link = $26,
			deezer_link = $27,
			other_link = $28,

			lyrics = $29,
			description = $30,
			notes = $31,

			updated_at = CURRENT_TIMESTAMP
		WHERE song_id = $31
	`

	_, err := DB.Exec(
		query,
		song.Title,
		song.TitleSlug,

		song.AlbumTitle,
		song.AlbumID,
		song.AlbumSlug,
		song.ArtistName,
		song.ArtistID,
		song.ArtistSlug,

		song.ArtworkID,
		song.ArtworkPath,
		song.ReleaseDate,
		song.Genre,

		song.RecordingBPM,
		song.LiveBPM,
		song.MusicalKey,
		song.Tuning,
		song.Capo,
		song.LengthSeconds,

		song.Status,
		song.Explicit,
		song.IsCover,

		song.SpotifyLink,
		song.AppleMusicLink,
		song.YouTubeLink,
		song.AmazonMusicLink,
		song.PandoraLink,
		song.DeezerLink,
		song.OtherLink,

		song.Lyrics,
		song.Description,
		song.Notes,

		song.SongID,
	)
	if err != nil {
		log.Println("   Unable to update song: ", err)
		return err
	}

	return nil
}

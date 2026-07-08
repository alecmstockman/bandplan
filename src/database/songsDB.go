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
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by TEXT,
		updated_by TEXT
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
	log.Println("- SongsTableCreateSong")

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
		time_signature,
		original_key,
		live_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,
		chords,
		chart_link,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		tidal_link,
		other_link,

		lyrics,
		description,
		notes,
		created_by,
		updated_by
	)
	VALUES (
		$1, $2,
		$3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20, $21, $22,
		$23, $24, $25, $26, $27,
		$28, $29, $30, $31, $32, $33, $34, $35,
		$36, $37, $38, $39, $40
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
		time_signature,
		original_key,
		live_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,
		chords,
		chart_link,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		tidal_link,
		other_link,

		lyrics,
		description,
		notes,
		created_at,
		updated_at,
		created_by,
		updated_by
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
		song.TimeSignature,
		song.OriginalKey,
		song.LiveKey,
		song.Tuning,
		song.Capo,
		song.LengthSeconds,

		song.Status,
		song.Explicit,
		song.IsCover,
		song.Chords,
		song.ChartLink,

		song.SpotifyLink,
		song.AppleMusicLink,
		song.YouTubeLink,
		song.AmazonMusicLink,
		song.PandoraLink,
		song.DeezerLink,
		song.TidalLink,
		song.OtherLink,

		song.Lyrics,
		song.Description,
		song.Notes,
		song.CreatedBy,
		song.UpdatedBy,
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
		&newSong.TimeSignature,
		&newSong.OriginalKey,
		&newSong.LiveKey,
		&newSong.Tuning,
		&newSong.Capo,
		&newSong.LengthSeconds,

		&newSong.Status,
		&newSong.Explicit,
		&newSong.IsCover,
		&newSong.Chords,
		&newSong.ChartLink,

		&newSong.SpotifyLink,
		&newSong.AppleMusicLink,
		&newSong.YouTubeLink,
		&newSong.AmazonMusicLink,
		&newSong.PandoraLink,
		&newSong.DeezerLink,
		&newSong.TidalLink,
		&newSong.OtherLink,

		&newSong.Lyrics,
		&newSong.Description,
		&newSong.Notes,
		&newSong.CreatedAt,
		&newSong.UpdatedAt,
		&newSong.CreatedBy,
		&newSong.UpdatedBy,
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
		time_signature,
		original_key,
		live_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,
		chords,
		chart_link,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		tidal_link,
		other_link,

		lyrics,
		description,
		notes,
		created_at,
		updated_at,
		created_by,
		updated_by
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
			&song.ArtistName,
			&song.ArtistID,
			&song.ArtistSlug,

			&song.ArtworkID,
			&song.ArtworkPath,
			&song.ReleaseDate,
			&song.Genre,

			&song.RecordingBPM,
			&song.LiveBPM,
			&song.TimeSignature,
			&song.OriginalKey,
			&song.LiveKey,
			&song.Tuning,
			&song.Capo,
			&song.LengthSeconds,

			&song.Status,
			&song.Explicit,
			&song.IsCover,
			&song.Chords,
			&song.ChartLink,

			&song.SpotifyLink,
			&song.AppleMusicLink,
			&song.YouTubeLink,
			&song.AmazonMusicLink,
			&song.PandoraLink,
			&song.DeezerLink,
			&song.TidalLink,
			&song.OtherLink,

			&song.Lyrics,
			&song.Description,
			&song.Notes,
			&song.CreatedAt,
			&song.UpdatedAt,
			&song.CreatedBy,
			&song.UpdatedBy,
		)
		if err != nil {
			log.Println("   Unable to get songs by band id: ", err)
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
	log.Println("- SongsTableSearchByBandID")

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
			time_signature,
			original_key,
			live_key,
			tuning,
			capo,
			length_seconds,

			status,
			explicit,
			is_cover,
			chords,
			chart_link,

			spotify_link,
			apple_music_link,
			youtube_link,
			amazon_music_link,
			pandora_link,
			deezer_link,
			tidal_link,
			other_link,

			lyrics,
			description,
			notes,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM songs
		WHERE band_id = $1
		AND (
			$2 = ''
			OR title ILIKE '%' || $2 || '%'
			OR album_title ILIKE '%' || $2 || '%'
			OR artist_name ILIKE '%' || $2 || '%'
			OR genre ILIKE '%' || $2 || '%'
			OR time_signature ILIKE '%' || $2 || '%'
			OR original_key ILIKE '%' || $2 || '%'
			OR live_key ILIKE '%' || $2 || '%'
			OR tuning ILIKE '%' || $2 || '%'
			OR chords ILIKE '%' || $2 || '%'
			OR recording_bpm::text ILIKE '%' || $2 || '%'
			OR live_bpm::text ILIKE '%' || $2 || '%'
			OR length_seconds::text ILIKE '%' || $2 || '%'
			OR lyrics ILIKE '%' || $2 || '%'
			OR notes ILIKE '%' || $2 || '%'
		)
		ORDER BY title ASC
	`, bandID, query)

	if err != nil {
		log.Println("   Unable to search songs by query: ", err)
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
			&song.ArtistName,
			&song.ArtistID,
			&song.ArtistSlug,

			&song.ArtworkID,
			&song.ArtworkPath,
			&song.ReleaseDate,
			&song.Genre,

			&song.RecordingBPM,
			&song.LiveBPM,
			&song.TimeSignature,
			&song.OriginalKey,
			&song.LiveKey,
			&song.Tuning,
			&song.Capo,
			&song.LengthSeconds,

			&song.Status,
			&song.Explicit,
			&song.IsCover,
			&song.Chords,
			&song.ChartLink,

			&song.SpotifyLink,
			&song.AppleMusicLink,
			&song.YouTubeLink,
			&song.AmazonMusicLink,
			&song.PandoraLink,
			&song.DeezerLink,
			&song.TidalLink,
			&song.OtherLink,

			&song.Lyrics,
			&song.Description,
			&song.Notes,
			&song.CreatedAt,
			&song.UpdatedAt,
			&song.CreatedBy,
			&song.UpdatedBy,
		)
		fmt.Println("song: ", song.Title)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	fmt.Println("rows.Err", rows.Err())
	return songs, rows.Err()
}

func SongsTableGetSongBySongID(songID string) (models.Song, error) {
	log.Println("- SongsTableGetSongBySongID")

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
		time_signature,
		original_key,
		live_key,
		tuning,
		capo,
		length_seconds,

		status,
		explicit,
		is_cover,
		chords,
		chart_link,

		spotify_link,
		apple_music_link,
		youtube_link,
		amazon_music_link,
		pandora_link,
		deezer_link,
		tidal_link,
		other_link,

		lyrics,
		description,
		notes,
		created_at,
		updated_at,
		created_by,
		updated_by
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
		&song.ArtistName,
		&song.ArtistID,
		&song.ArtistSlug,

		&song.ArtworkID,
		&song.ArtworkPath,
		&song.ReleaseDate,
		&song.Genre,

		&song.RecordingBPM,
		&song.LiveBPM,
		&song.TimeSignature,
		&song.OriginalKey,
		&song.LiveKey,
		&song.Tuning,
		&song.Capo,
		&song.LengthSeconds,

		&song.Status,
		&song.Explicit,
		&song.IsCover,
		&song.Chords,
		&song.ChartLink,

		&song.SpotifyLink,
		&song.AppleMusicLink,
		&song.YouTubeLink,
		&song.AmazonMusicLink,
		&song.PandoraLink,
		&song.DeezerLink,
		&song.TidalLink,
		&song.OtherLink,

		&song.Lyrics,
		&song.Description,
		&song.Notes,
		&song.CreatedAt,
		&song.UpdatedAt,
		&song.CreatedBy,
		&song.UpdatedBy,
	)
	if err != nil {
		log.Println("   Unable to get song from songs db: ", err)
		return models.Song{}, nil
	}

	return song, nil
}

func SongsTableUpdateSong(song models.Song) error {
	log.Println("- SongsTableUpdateSongWithArt")

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
			time_signature = $15,
			original_key = $16,
			live_key = $17,
			tuning = $18,
			capo = $19,
			length_seconds = $20,

			status = $21,
			explicit = $22,
			is_cover = $23,
			chords = $24,
			chart_link = $25,

			spotify_link = $26,
			apple_music_link = $27,
			youtube_link = $28,
			amazon_music_link = $29,
			pandora_link = $30,
			deezer_link = $31,
			tidal_link = $32,
			other_link = $33,

			lyrics = $34,
			description = $35,
			notes = $36,
			updated_by = $37,

			updated_at = CURRENT_TIMESTAMP
		WHERE song_id = $38
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
		song.TimeSignature,
		song.OriginalKey,
		song.LiveKey,
		song.Tuning,
		song.Capo,
		song.LengthSeconds,

		song.Status,
		song.Explicit,
		song.IsCover,
		song.Chords,
		song.ChartLink,

		song.SpotifyLink,
		song.AppleMusicLink,
		song.YouTubeLink,
		song.AmazonMusicLink,
		song.PandoraLink,
		song.DeezerLink,
		song.TidalLink,
		song.OtherLink,

		song.Lyrics,
		song.Description,
		song.Notes,
		song.UpdatedBy,

		song.SongID,
	)
	if err != nil {
		log.Println("   Unable to update song: ", err)
		return err
	}

	return nil
}

func SongsTableUpdateSongWithoutArt(song models.Song) error {
	log.Println("- SongsTableUpdateSongWithArt")

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

			release_date = $9,
			genre = $10,

			recording_bpm = $11,
			live_bpm = $12,
			time_signature = $13,
			original_key = $14,
			live_key = $15,
			tuning = $16,
			capo = $17,
			length_seconds = $18,

			status = $19,
			explicit = $20,
			is_cover = $21,
			chords = $22,
			chart_link = $23,

			spotify_link = $24,
			apple_music_link = $25,
			youtube_link = $26,
			amazon_music_link = $27,
			pandora_link = $28,
			deezer_link = $29,
			tidal_link = $30,
			other_link = $31,

			lyrics = $32,
			description = $33,
			notes = $34,
			updated_by = $35,

			updated_at = CURRENT_TIMESTAMP
		WHERE song_id = $36
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

		song.ReleaseDate,
		song.Genre,

		song.RecordingBPM,
		song.LiveBPM,
		song.TimeSignature,
		song.OriginalKey,
		song.LiveKey,
		song.Tuning,
		song.Capo,
		song.LengthSeconds,

		song.Status,
		song.Explicit,
		song.IsCover,
		song.Chords,
		song.ChartLink,

		song.SpotifyLink,
		song.AppleMusicLink,
		song.YouTubeLink,
		song.AmazonMusicLink,
		song.PandoraLink,
		song.DeezerLink,
		song.TidalLink,
		song.OtherLink,

		song.Lyrics,
		song.Description,
		song.Notes,
		song.UpdatedBy,

		song.SongID,
	)
	if err != nil {
		log.Println("   Unable to update song: ", err)
		return err
	}

	return nil
}

func SongsTableDeleteSongByID(songID string) error {
	log.Println("- SongsTableDeleteSongByID")

	query := `
	DELETE FROM songs
	WHERE song_id = $1
	`
	_, err := DB.Exec(query, songID)
	if err != nil {
		log.Println("   Unable to delete song from songs db: ", err)
		return err
	}
	return nil
}

func SongsTableGetImageIDAndPathBySongID(songID string) (string, string, error) {
	log.Println("- SongsTableGetImageIDBySongID")

	query := `
	SELECT artwork_id, artwork_path
	FROM songs
	WHERE song_id = $1
	`

	var artworkID string
	var artworkPath string

	err := DB.QueryRow(query, songID).Scan(&artworkID, &artworkPath)
	if err != nil {
		log.Println("   Unable to get artwork ID or Path from songs table: ", err)
		return "", "", nil
	}

	return artworkID, artworkPath, nil
}

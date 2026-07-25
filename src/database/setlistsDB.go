package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func SetlistsTableGetSetlistsByBandID(bandID string) ([]models.Setlist, error) {
	log.Println("- SetlistsTableGetSetlistByBandID")

	query := `
	SELECT *
	FROM setlists
	WHERE band_id = $1
	`

	rows, err := DB.Query(query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	setlists := []models.Setlist{}

	for rows.Next() {
		var setlist models.Setlist

		err := rows.Scan(
			&setlist.ID,
			&setlist.SetlistID,
			&setlist.BandID,
			&setlist.Name,
			&setlist.Notes,
			&setlist.ArtworkPath,
			&setlist.CreatedAt,
			&setlist.CreatedBy,
			&setlist.UpdatedAt,
			&setlist.UpdatedBy,
		)
		if err != nil {
			log.Println("   Unable to get setlist by band id: ", err)
			return nil, err
		}
		setlists = append(setlists, setlist)
	}

	return setlists, nil
}

func SetlistsTableCreateSetlist(setlist models.Setlist) error {
	log.Println("- SetlistsTableCreateSetlist")

	setlistID := uuid.New().String()

	// fmt.Println(" Setlist: ", setlist)
	// fmt.Println(" user: ", setlist.UpdatedBy)

	query := `
		INSERT INTO setlists(
			setlist_id,
			band_id,
			name,
			notes,
			artwork_path,
			created_by,
			updated_by
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := DB.Exec(
		query,
		setlistID,
		setlist.BandID,
		setlist.Name,
		setlist.Notes,
		setlist.ArtworkPath,
		setlist.CreatedBy,
		setlist.UpdatedBy,
	)

	if err != nil {
		log.Println("   Unable to create setlist in database: ", err)
		return err
	}
	return nil
}

func SetlistsTableGetSetlistByID(setlistID string) (models.Setlist, error) {
	log.Println("- SetlistsTableGetSetlistByID")

	setlistQuery := `
		SELECT 
			id,
			setlist_id,
			band_id,
			name,
			notes,
			artwork_path,
			created_at,
			created_by,
			updated_at,
			updated_by

		FROM setlists
		WHERE
			setlist_id = $1
	`

	setlist := models.Setlist{}

	err := DB.QueryRow(setlistQuery, setlistID).Scan(
		&setlist.ID,
		&setlist.SetlistID,
		&setlist.BandID,
		&setlist.Name,
		&setlist.Notes,
		&setlist.ArtworkPath,
		&setlist.CreatedAt,
		&setlist.CreatedBy,
		&setlist.UpdatedAt,
		&setlist.UpdatedBy,
	)

	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return models.Setlist{}, err
	}

	songQuery := `
		SELECT 
			ss.id,
			ss.setlist_id,
			ss.song_id,
			ss.position,
			ss.created_at,
			ss.created_by,
			ss.updated_at,
			ss.updated_by,
			
			s.id,
			s.song_id,
			s.band_id,
			s.title,
			s.title_slug,
			s.album_title,
			s.album_id,
			s.album_slug,
			s.artist_name,
			s.artist_id,
			s.artist_slug,
			s.artwork_id,
			s.artwork_path,
			s.release_date,
			s.genre,
			s.recording_bpm,
			s.live_bpm,
			s.time_signature,
			s.original_key,
			s.live_key,
			s.tuning,
			s.capo,
			s.length_seconds,
			s.status,
			s.explicit,
			s.is_cover,
			s.chords,
			s.chart_link,
			s.spotify_link,
			s.apple_music_link,
			s.youtube_link,
			s.amazon_music_link,
			s.pandora_link,
			s.deezer_link,
			s.tidal_link,
			s.other_link,
			s.lyrics,
			s.description,
			s.notes,
			s.created_at,
			s.created_by,
			s.updated_at,
			s.updated_by

		FROM setlist_songs ss
		LEFT JOIN songs s 
			ON ss.song_id = s.song_id
		WHERE ss.setlist_id = $1
		ORDER BY ss.position
	`

	rows, err := DB.Query(songQuery, setlistID)
	if err != nil {
		log.Println("   Unable to get song from table by ID: ", err)
		return models.Setlist{}, err
	}
	defer rows.Close()

	for rows.Next() {

		setlistSong := models.SetlistItem{}
		song := models.Song{}

		err := rows.Scan(
			&setlistSong.ID,
			&setlistSong.SetlistID,
			&setlistSong.SongID,
			&setlistSong.Position,
			&setlistSong.CreatedAt,
			&setlistSong.CreatedBy,
			&setlistSong.UpdatedAt,
			&setlistSong.UpdatedBy,

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
			&song.CreatedBy,
			&song.UpdatedAt,
			&song.UpdatedBy,
		)
		if err != nil {
			log.Println(" Unable to get setlist: ", err)
			return models.Setlist{}, err
		}

		setlistSong.Song = song

		fmt.Println("\n=============================")
		fmt.Println(song.ID)
		fmt.Println(song.SongID)
		fmt.Println(song.Title)
		fmt.Println("=============================")
		fmt.Println(setlistSong.Song.Title)
		fmt.Println("=============================\n")

		setlist.Songs = append(setlist.Songs, setlistSong)
	}

	for _, song := range setlist.Songs {
		fmt.Println("song: ", song.Song.Title)
	}

	return setlist, nil
}

func SetlistsTableDeleteSetlist(setlistID string) error {
	log.Println("- SetlistsTableDeleteSetlist")

	query := `
	DELETE FROM setlists
	WHERE setlist_id = $1
	`

	_, err := DB.Exec(query, setlistID)
	if err != nil {
		log.Printf("\n   Unable to delete setlist: %s, err: ", setlistID, err)
		return err
	}

	return nil
}

func SetlistItemsTableSaveSong(songID string, userID string, setlistID string) (models.SetlistItem, error) {
	log.Println("- SetlistItemsTableSaveSong")

	query := `
		INSERT INTO setlist_songs(
			setlist_id,
			song_id,
			position,
			created_by,
			updated_by
		) 
		VALUES (
			$1,
			$2,
			COALESCE(
				(
					SELECT MAX(position) + 1
					FROM setlist_songs
					WHERE setlist_id = $1
				),
				1
			),
			$3,
			$4
		)
		RETURNING
			id, 
			setlist_id,
			song_id,
			position,
			created_at,
			created_by,
			updated_at,
			updated_by
	`

	song := models.SetlistItem{}

	err := DB.QueryRow(
		query,
		setlistID,
		songID,
		userID,
		userID,
	).Scan(
		&song.ID,
		&song.SetlistID,
		&song.SongID,
		&song.Position,
		&song.CreatedAt,
		&song.CreatedBy,
		&song.UpdatedAt,
		&song.UpdatedBy,
	)

	if err != nil {
		log.Printf("\n   Unable to save %s in setlist_songs table: %v", songID, err)
		return models.SetlistItem{}, err
	}

	return song, nil
}

func SetlistItemsTableDeleteSong(songID string, setlistID string) error {
	log.Println("- SetlistItemsTableDeleteSong")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var deletedPosition int

	deleteQuery := `
		DELETE FROM setlist_songs
		WHERE song_id = $1
			AND setlist_id = $2
		RETURNING position
	`

	err = tx.QueryRow(deleteQuery, songID, setlistID).Scan(&deletedPosition)
	if err != nil {
		return err
	}

	UpdateQuery := `
		UPDATE setlist_songs
		SET
			position = position - 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE setlist_id = $1
			AND position > $2
	`

	_, err = tx.Exec(
		UpdateQuery,
		setlistID,
		deletedPosition,
	)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func SetlistsSongsTableGetAllSongsBySetlistID(setlistID string) ([]models.SetlistItem, error) {
	log.Println("- SetlistsSongsTableGetAllSongsByBandID")

	query := `
		SELECT
			id, 
			setlist_id,
			song_id,
			position,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM setlist_songs
		WHERE setlist_id = $1
	`

	rows, err := DB.Query(
		query,
		setlistID,
	)
	if err != nil {
		log.Println("   Unable to query setlist_songs: ", err)
		return []models.SetlistItem{}, err
	}

	defer rows.Close()

	songsList := []models.SetlistItem{}

	for rows.Next() {
		song := models.SetlistItem{}

		err := rows.Scan(
			&song.ID,
			&song.SetlistID,
			&song.SongID,
			&song.Position,
			&song.CreatedAt,
			&song.CreatedBy,
			&song.UpdatedAt,
			&song.UpdatedBy,
		)
		if err != nil {
			log.Printf("\n   Unable to get songs from setlist_songs with setlistID :", setlistID, err)
			return []models.SetlistItem{}, err
		}
		songsList = append(songsList, song)
	}

	return songsList, nil
}

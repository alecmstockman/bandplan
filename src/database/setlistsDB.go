package database

import (
	"bandplan/src/models"
	"log"

	"github.com/google/uuid"
)

func SetlistsTableGetSetlistsByBandID(bandID string) ([]models.Setlist, error) {
	log.Println("- SetlistsTableGetSetlistByBandID")

	query := `
	SELECT
		id,
		setlist_id,
		band_id,
		name,
		COALESCE(slug, ''),
		explicit,
		COALESCE(notes, ''),
		COALESCE(image_id, ''),
		COALESCE(artwork_path, ''),
		created_at,
		created_by,
		updated_at,
		updated_by
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
			&setlist.Slug,
			&setlist.Explicit,
			&setlist.Notes,
			&setlist.ArtworkID,
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

func SetlistsTableGetSetlistSummariessByBandID(bandID string) ([]models.SetlistSummary, error) {
	log.Println("- SetlistsTableGetSetlistItemsByBandID")

	query := `
		SELECT 
			s.setlist_id,
			s.name,
			s.slug,

			COUNT(ss.song_id) AS song_count,
			COALESCE(SUM(so.length_seconds), 0) AS total_length_seconds,

			COALESCE(s.image_id, ''),
			COALESCE(s.artwork_path, ''),
			s.created_at,
			s.created_by,
			s.updated_at,
			s.updated_by
			FROM 
				setlists s
			LEFT JOIN setlist_items ss
				ON s.setlist_id = ss.setlist_id

			LEFT JOIN songs so
				ON ss.song_id = so.song_id

			WHERE 
				s.band_id = $1
			GROUP BY
				s.setlist_id,
				s.name,
				s.slug,
				s.image_id,
				s.artwork_path,
				s.created_at,
				s.created_by,
				s.updated_at,
				s.updated_by
			ORDER BY
				s.name ASC;
	`
	rows, err := DB.Query(
		query,
		bandID,
	)
	if err != nil {
		log.Println("   Unable to get setlist summaries from db: ", err)
		return nil, err
	}

	defer rows.Close()

	setlists := []models.SetlistSummary{}

	for rows.Next() {
		setlist := models.SetlistSummary{}

		err := rows.Scan(
			&setlist.SetlistID,
			&setlist.Name,
			&setlist.Slug,
			&setlist.SongCount,
			&setlist.Length,
			&setlist.ArtworkID,
			&setlist.ArtworkPath,
			&setlist.CreatedAt,
			&setlist.CreatedBy,
			&setlist.UpdatedAt,
			&setlist.UpdatedBy,
		)
		if err != nil {
			log.Println("   Unable to scan setlist summaries from db: ", err)
			return []models.SetlistSummary{}, err
		}
		setlist.BandID = bandID
		setlists = append(setlists, setlist)
	}

	return setlists, nil
}

func SetlistsTableCreateSetlist(setlist models.Setlist) error {
	log.Println("- SetlistsTableCreateSetlist")

	setlistID := uuid.New().String()

	query := `
		INSERT INTO setlists(
			setlist_id,
			band_id,
			name,
			slug,
			explicit,
			notes,
			image_id,
			artwork_path,
			created_by,
			updated_by
		)
		VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, $10
		)
	`

	_, err := DB.Exec(
		query,
		setlistID,
		setlist.BandID,
		setlist.Name,
		setlist.Slug,
		setlist.Explicit,
		setlist.Notes,
		setlist.ArtworkID,
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

func SetlistsTableUpdateSetlist(setlist models.Setlist) error {
	log.Println("- SetlistsTableUpdateSetlist")
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
			slug,
			explicit,
			notes,
			image_id,
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
		&setlist.Slug,
		&setlist.Explicit,
		&setlist.Notes,
		&setlist.ArtworkID,
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

		FROM setlist_items ss
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

		setlist.Songs = append(setlist.Songs, setlistSong)
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
		log.Printf("\n   Unable to delete setlist: %s, err: %v", setlistID, err)
		return err
	}

	return nil
}

func SetlistItemsTableSaveItem(itemType models.SetlistItemType, itemID string, userID string, setlistID string) (models.SetlistItem, error) {
	log.Println("- SetlistItemsTableSaveItem")

	var itemTypeString string

	switch itemType {
	case models.SetlistItemSong:
		itemTypeString = "song"
	case models.SetlistItemTransition:
		itemTypeString = "transition"
	case models.SetlistItemBreak:
		itemTypeString = "break"
	}

	query := `
		INSERT INTO setlist_items(
			setlist_id,
			song_id,
			position,
			created_by,
			updated_by,
			item_type
		) 
		VALUES (
			$1,
			$2,
			COALESCE(
				(
					SELECT MAX(position) + 1
					FROM setlist_items
					WHERE setlist_id = $1
				),
				1
			),
			$3,
			$4,
			$5
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

	item := models.SetlistItem{}

	err := DB.QueryRow(
		query,
		setlistID,
		itemID,
		userID,
		userID,
		itemTypeString,
	).Scan(
		&item.ID,
		&item.SetlistID,
		&item.SongID,
		&item.Position,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.UpdatedAt,
		&item.UpdatedBy,
	)

	if err != nil {
		log.Printf("\n   Unable to save %s in setlist_items table: %v", itemID, err)
		return models.SetlistItem{}, err
	}

	return item, nil
}

func SetlistItemsTableDeleteSong(songID string, position int, setlistID string) error {
	log.Println("- SetlistItemsTableDeleteSong")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var deletedPosition int

	deleteQuery := `
		DELETE FROM setlist_items
		WHERE song_id = $1
			AND position = $2
			AND setlist_id = $3
		RETURNING position
	`

	err = tx.QueryRow(
		deleteQuery,
		songID,
		position,
		setlistID,
	).Scan(&deletedPosition)

	if err != nil {
		return err
	}

	UpdateQuery := `
		UPDATE setlist_items
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
		FROM setlist_items
		WHERE setlist_id = $1
	`

	rows, err := DB.Query(
		query,
		setlistID,
	)
	if err != nil {
		log.Println("   Unable to query setlist_items: ", err)
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
			log.Println("\n   Unable to get songs from setlist_items with setlistID :", setlistID, err)
			return []models.SetlistItem{}, err
		}
		songsList = append(songsList, song)
	}

	return songsList, nil
}

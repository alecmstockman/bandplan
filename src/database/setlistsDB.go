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
		WHERE setlist_id = $1
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

	itemQuery := `
		SELECT
			id,
			setlist_id,
			COALESCE(
				song_id,
				transition_id,
				break_id
			) AS item_id,
			item_type,
			position,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM setlist_items
		WHERE setlist_id = $1
		ORDER BY position
	`

	rows, err := DB.Query(itemQuery, setlistID)
	if err != nil {
		log.Println("   Unable to get setlist items: ", err)
		return models.Setlist{}, err
	}
	defer rows.Close()

	for rows.Next() {
		item := models.SetlistItem{}

		err := rows.Scan(
			&item.ID,
			&item.SetlistID,
			&item.ItemID,
			&item.ItemType,
			&item.Position,
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
		)
		if err != nil {
			log.Println("   Unable to scan setlist item: ", err)
			return models.Setlist{}, err
		}

		switch item.ItemType {

		case models.SetlistItemSong:
			songQuery := `
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
					created_by,
					updated_at,
					updated_by
				FROM songs
				WHERE song_id = $1
			`

			song := models.Song{}

			err := DB.QueryRow(songQuery, item.ItemID).Scan(
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
				log.Println("   Unable to get song for setlist item: ", err)
				return models.Setlist{}, err
			}

			item.Song = song

		case models.SetlistItemTransition:
			// Load transition here once you have the
			// Transition model/query set up.

		case models.SetlistItemBreak:
			// Load break here once you have the
			// Break model/query set up.

		default:
			return models.Setlist{}, fmt.Errorf(
				"unknown setlist item type: %s",
				item.ItemType,
			)
		}

		setlist.Songs = append(setlist.Songs, item)
	}

	if err := rows.Err(); err != nil {
		log.Println("   Error iterating setlist items: ", err)
		return models.Setlist{}, err
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
	var query string

	item := models.SetlistItem{}

	switch itemType {
	case models.SetlistItemSong:
		itemTypeString = "song"
		query = `
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
	case models.SetlistItemTransition:
		itemTypeString = "transition"
		query = `
			INSERT INTO setlist_items(
				setlist_id,
				transition_id,
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
				transition_id,
				position,
				created_at,
				created_by,
				updated_at,
				updated_by
		`
	case models.SetlistItemBreak:
		itemTypeString = "break"
		query = `
			INSERT INTO setlist_items(
				setlist_id,
				break_id,
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
				break_id,
				position,
				created_at,
				created_by,
				updated_at,
				updated_by
		`
	}

	var err error

	row := DB.QueryRow(
		query,
		setlistID,
		itemID,
		userID,
		userID,
		itemTypeString,
	)

	switch itemType {
	case models.SetlistItemSong:
		err = row.Scan(
			&item.ID,
			&item.SetlistID,
			&item.ItemID,
			&item.Position,
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
		)

	case models.SetlistItemTransition:
		err = row.Scan(
			&item.ID,
			&item.SetlistID,
			&item.ItemID,
			&item.Position,
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
		)

	case models.SetlistItemBreak:
		err = row.Scan(
			&item.ID,
			&item.SetlistID,
			&item.ItemID,
			&item.Position,
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
		)
	default:
		return models.SetlistItem{}, fmt.Errorf(
			"invalid setlist item type: %s",
			itemType,
		)
	}

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
			&song.ItemID,
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

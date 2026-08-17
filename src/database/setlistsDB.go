package database

import (
	"bandplan/src/models"
	"errors"
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
			pause_after_seconds,
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
			&item.PauseAfterSeconds,
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

			item.Song = &song

		case models.SetlistItemTransition:
			transitionQuery := `
				SELECT
					id,
					transition_id,
					band_id,
					title,
					title_slug,
					length_seconds,
					bpm,
					time_signature,
					musical_key,
					tuning,
					capo,
					explicit,
					chords,
					chart_link,
					lyrics,
					notes,
					link_one,
					link_two,
					link_three,
					created_at,
					created_by,
					updated_at,
					updated_by

				FROM transitions
				WHERE transition_id = $1
			`
			transition := models.Transition{}

			err := DB.QueryRow(transitionQuery, item.ItemID).Scan(
				&transition.ID,
				&transition.TransitionID,
				&transition.BandID,
				&transition.Title,
				&transition.Slug,
				&transition.LengthSeconds,
				&transition.BPM,
				&transition.TimeSignature,
				&transition.Key,
				&transition.Tuning,
				&transition.Capo,
				&transition.Explicit,
				&transition.Chords,
				&transition.ChartLink,
				&transition.Lyrics,
				&transition.Notes,
				&transition.LinkOne,
				&transition.LinkTwo,
				&transition.LinkThree,
				&transition.CreatedAt,
				&transition.CreatedBy,
				&transition.UpdatedAt,
				&transition.UpdatedBy,
			)
			if err != nil {
				log.Println("   Unable to get transition for setlist item: ", err)
				return models.Setlist{}, err
			}

			item.Transition = &transition

		case models.SetlistItemBreak:
			breakQuery := `
				SELECT 
					id,
					break_id,
					band_id,
					title,
					title_slug,
					notes,
					length_seconds,
					link_one,
					link_two,
					created_at,
					created_by,
					updated_at,
					updated_by

				FROM breaks
				WHERE break_id = $1
			`
			breakItem := models.Break{}
			err := DB.QueryRow(breakQuery, item.ItemID).Scan(
				&breakItem.ID,
				&breakItem.BreakID,
				&breakItem.BandID,
				&breakItem.Title,
				&breakItem.Slug,
				&breakItem.Notes,
				&breakItem.LengthSeconds,
				&breakItem.LinkOne,
				&breakItem.LinkTwo,
				&breakItem.CreatedAt,
				&breakItem.CreatedBy,
				&breakItem.UpdatedAt,
				&breakItem.UpdatedBy,
			)
			if err != nil {
				log.Println("   Unable to get break for setlist: ", err)
				return models.Setlist{}, err
			}

			item.Break = &breakItem
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

	var itemColumn string

	switch itemType {
	case models.SetlistItemSong:
		itemColumn = "song_id"

	case models.SetlistItemTransition:
		itemColumn = "transition_id"

	case models.SetlistItemBreak:
		itemColumn = "break_id"

	default:
		return models.SetlistItem{}, fmt.Errorf(
			"invalid setlist item type: %s",
			itemType,
		)
	}

	fmt.Println("==== item column: ", itemColumn)

	query := fmt.Sprintf(`
		INSERT INTO setlist_items (
			setlist_id,
			%s,
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
			%s,
			position,
			created_at,
			created_by,
			updated_at,
			updated_by
	`, itemColumn, itemColumn)

	item := models.SetlistItem{
		ItemType: itemType,
	}

	err := DB.QueryRow(
		query,
		setlistID,
		itemID,
		userID,
		userID,
		string(itemType),
	).Scan(
		&item.ID,
		&item.SetlistID,
		&item.ItemID,
		&item.Position,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.UpdatedAt,
		&item.UpdatedBy,
	)

	if err != nil {
		log.Printf(
			"   Unable to save %s in setlist_items table: %v",
			itemID,
			err,
		)

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

func SetlistItemsTableDeleteTransition(transitionID string, position int, setlistID string) error {
	log.Println("- SetlistItemsTableDeleteTransition")

	fmt.Println("transitionID: ", transitionID)
	fmt.Println("position: ", position)
	fmt.Println("setlistID: ", setlistID)

	tx, err := DB.Begin()
	if err != nil {
		log.Println("   Unable to begin database transaction to delete transition: ", err)
		return err
	}
	defer tx.Rollback()

	var deletedPosition int

	deleteQuery := `
		DELETE FROM setlist_items
		WHERE transition_id = $1
			AND position = $2
			AND setlist_id = $3
		RETURNING position
	`

	err = tx.QueryRow(
		deleteQuery,
		transitionID,
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
		fmt.Println("   err: ", err)
		return err
	}

	return tx.Commit()

}

func SetlistItemsTableDeleteBreak(breakID string, position int, setlistID string) error {
	log.Println("- SetlistItemsTableDeleteTransition")

	fmt.Println("transitionID: ", breakID)
	fmt.Println("position: ", position)
	fmt.Println("setlistID: ", setlistID)

	tx, err := DB.Begin()
	if err != nil {
		log.Println("   Unable to begin database transaction to delete break: ", err)
		return err
	}
	defer tx.Rollback()

	var deletedPosition int

	deleteQuery := `
		DELETE FROM setlist_items
		WHERE break_id = $1
			AND position = $2
			AND setlist_id = $3
		RETURNING position
	`

	err = tx.QueryRow(
		deleteQuery,
		breakID,
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
		fmt.Println("   err: ", err)
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

func SetlistItemsUpdateOrder(setlistID string, newOrder []models.ReorderItem) error {
	log.Println("- SetlistItemsUpdateOrder")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateQuery := `
	UPDATE setlist_items
	SET 
		position = $1
	WHERE id = $2 
		AND setlist_id = $3
	`

	for i, item := range newOrder {
		_, err := tx.Exec(
			updateQuery,
			-(i + 1),
			item.ItemID,
			setlistID,
		)
		if err != nil {
			log.Println("   Unable to assign temp order while updatign order: ", err)
			return err
		}

	}
	for _, item := range newOrder {
		_, err = tx.Exec(
			updateQuery,
			item.Position,
			item.ItemID,
			setlistID,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("   Unable to save new order to db: ", err)
		return err
	}

	return nil
}

func SetlistItemsGetItem(setlistID string, itemType models.SetlistItemType, itemID string) (models.SetlistItem, error) {
	fmt.Println("-------------------------")
	log.Println("- SetlistItemsGetItem")

	var query string

	fmt.Println("setlistID: ", setlistID)
	fmt.Println("itemType: ", itemType)
	fmt.Println("itemID: ", itemID)

	if itemType == "song" {
		query = `
		SELECT 
			id,
			setlist_id,
			item_type,
			song_id,
			position,
			pause_after_seconds,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM setlist_items
		WHERE setlist_id = $1
			AND song_id = $2
		`
	}
	if itemType == "transition" {
		query = `
			SELECT 
				id,
				setlist_id,
				item_type,
				transition_id,
				position,
				pause_after_seconds,
				created_at,
				created_by,
				updated_at,
				updated_by
			FROM setlist_items
			WHERE setlist_id = $1
				AND transition_id = $2
			`
	}
	if itemType == "break" {
		query = `
			SELECT
				id,
				break_id,
				band_id,
				title,
				title_slug,
				notes,
				length_seconds,
				link_one,
				link_two,
				created_at, 
				created_by,
				updated_at,
				updated_by
			FROM setlist_items
			WHERE setlist = $1
				AND break_id = $2
			`
	}

	fmt.Println("query: ", query)

	var item models.SetlistItem

	err := DB.QueryRow(
		query,
		setlistID,
		itemID,
	).Scan(
		&item.ID,
		&item.SetlistID,
		&item.ItemType,
		&item.ItemID,
		&item.Position,
		&item.PauseAfterSeconds,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.UpdatedAt,
		&item.UpdatedBy,
	)

	if err != nil {
		log.Println("   Unable to get setlistItem from setlist_items: ", err)
		return models.SetlistItem{}, err
	}

	return item, nil
}

func SetlistItemsUpdateItem(setlistID string, itemType models.SetlistItemType, itemID string, pauseAfter int) error {
	log.Println("- SetlistItemsUpdateItem")

	var query string

	if itemType == "song" {
		query = `
		UPDATE setlist_items
		SET
			pause_after_seconds = $1
		WHERE setlist_id = $2
			AND song_id = $3
		`
	} else if itemType == "transition" {
		query = `
		UPDATE setlist_items
		SET
			pause_after_seconds = $1
		WHERE setlist_id = $2
			AND	transition_id = $3
		`
	} else if itemType == "break" {
		query = `
		UPDATE setlist_items
		SET
			pause_after_seconds = $1
		WHERE setlist_id = $2
			AND	break_id = $3
		`
	} else {
		log.Println("   Unable to update setlist item")
		return errors.New("Unable to update setlist item")
	}

	_, err := DB.Exec(
		query,
		pauseAfter,
		setlistID,
		itemID,
	)
	if err != nil {
		log.Println("   Unable to update setlistitem: ", err)
		return err
	}

	return nil

}

package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func SetlistsTableGetSetlistsByBandIDAndUserID(bandID string, userID string) ([]models.Setlist, error) {
	log.Println("- SetlistsTableGetSetlistsByBandIDAndUserID")

	query := `
		SELECT
			s.id,
			s.setlist_id,
			s.band_id,
			s.name,
			COALESCE(s.slug, ''),
			s.explicit,
			COALESCE(s.notes, ''),
			COALESCE(s.image_id, ''),
			COALESCE(s.artwork_path, ''),
			s.created_at,
			s.created_by,
			s.updated_at,
			s.updated_by
		FROM setlists s
		WHERE s.band_id = $1
		AND EXISTS (
			SELECT 1
			FROM band_members bm
			WHERE bm.band_id = s.band_id
			AND bm.user_id = $2
	)`

	rows, err := DB.Query(query, bandID, userID)
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
			return nil, err
		}

		setlists = append(setlists, setlist)
	}

	return setlists, rows.Err()
}

func SetlistsTableGetSetlistSummariesByBandIDAndUserID(bandID string, userID string) ([]models.SetlistSummary, error) {
	log.Println("- SetlistsTableGetSetlistSummariesByBandIDAndUserID")

	query := `
		SELECT 
			s.setlist_id,
			s.name,
			s.slug,

			COUNT(ss.song_id) AS song_count,

			COALESCE(
				SUM(
					COALESCE(so.length_seconds, 0) +
					COALESCE(tr.length_seconds, 0) + 
					COALESCE(br.length_seconds, 0) + 
					COALESCE(ss.pause_after_seconds, 0)
				),
				0
			) AS total_length_seconds,

			COALESCE(s.image_id, ''),
			COALESCE(s.artwork_path, ''),
			s.created_at,
			s.created_by,
			s.updated_at,
			s.updated_by

		FROM setlists s

		LEFT JOIN setlist_items ss
			ON s.setlist_id = ss.setlist_id

		LEFT JOIN songs so
			ON ss.song_id = so.song_id

		LEFT JOIN transitions tr
			ON ss.transition_id = tr.transition_id

		LEFT JOIN breaks br
			ON ss.break_id = br.break_id

		WHERE 
			s.band_id = $1
			AND EXISTS (
				SELECT 1
				FROM band_members bm
				WHERE bm.band_id = s.band_id
					AND bm.user_id = $2
			)

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
		userID,
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
			return nil, err
		}

		setlist.BandID = bandID

		setlists = append(setlists, setlist)
	}

	if err := rows.Err(); err != nil {
		log.Println("   Error iterating setlist summaries: ", err)
		return nil, err
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

func SetlistsTableUpdateSetlist(setlist models.Setlist, userID string) (bool, error) {
	log.Println("- SetlistsTableUpdateSetlist")

	query := `
		UPDATE setlists
		SET
			name = $1,
			slug = $2,
			explicit = $3,
			notes = $4,
			image_id = $5,
			artwork_path = $6,
			updated_at = NOW(),
			updated_by = $7
		WHERE setlist_id = $8
			AND band_id = $9
			AND EXISTS (
				SELECT 1
				FROM band_members bm
				WHERE bm.band_id = setlists.band_id
					AND bm.user_id = $10
			)
	`

	result, err := DB.Exec(
		query,
		setlist.Name,
		setlist.Slug,
		setlist.Explicit,
		setlist.Notes,
		setlist.ArtworkID,
		setlist.ArtworkPath,
		setlist.UpdatedBy,
		setlist.SetlistID,
		setlist.BandID,
		userID,
	)
	if err != nil {
		log.Println("   Unable to update setlist in database: ", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("   Unable to confirm setlist update: ", err)
		return false, err
	}

	return rowsAffected == 1, nil
}

func SetlistsTableGetSetlistByIDAndUserID(setlistID string, userID string) (models.Setlist, error) {
	log.Println("- SetlistsTableGetSetlistByIDAndUserID")

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
			AND EXISTS (
				SELECT 1
				FROM band_members bm
				WHERE bm.band_id = setlists.band_id
					AND bm.user_id = $2)
	`

	setlist := models.Setlist{}

	err := DB.QueryRow(setlistQuery, setlistID, userID).Scan(
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

func SetlistsTableUpdateNotes(setlistID string, newNotes string) error {
	log.Println("- SetlistTableUpdateNotes")

	query := `
		UPDATE setlists
		SET notes = $1
		WHERE setlist_id = $2
	`

	_, err := DB.Exec(query, newNotes, setlistID)
	if err != nil {
		log.Println("   Unable to save notes to DB: ", err)
		return err
	}
	return nil
}

func SetlistsTableSearchSetlistByBandIDAndUserID(bandID string, userID string, query string) ([]models.Setlist, error) {
	log.Println("- SetlistsTableSearchSetlistByID")

	rows, err := DB.Query(`
		SELECT
			id,
			setlist_id,
			band_id,
			name,
			notes,
			artwork_path,
			slug,
			explicit,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM setlists
		WHERE band_id = $1
			AND EXISTS (
				SELECT 1
				FROM band_members bm
				WHERE bm.band_id = s.band_id
					AND bm.user_id = $2
			)
			AND (
				$2 = ''
				OR name ILIKE '%' || $2 || '%'
				OR notes ILIKE '%' || $2 || '%'
			)
		ORDER BY name ASC
	`, bandID, query)

	if err != nil {
		log.Println("   Unable to search songs by query: ", err)
		return nil, err
	}
	defer rows.Close()

	var setlists []models.Setlist

	for rows.Next() {
		var setlist models.Setlist

		err := rows.Scan(
			&setlist.ID,
		)
		if err != nil {
			return nil, err
		}
		setlists = append(setlists, setlist)
	}

	return setlists, nil
}

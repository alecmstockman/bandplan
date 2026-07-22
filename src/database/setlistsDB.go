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

	fmt.Println(" Setlist: ", setlist)
	fmt.Println(" user: ", setlist.UpdatedBy)

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

func SetlistSongsTableSaveSong(songID string, userID string, setlistID string) (models.SetlistSong, error) {
	log.Println("- SetlistSonsTableSaveSetlist")

	query := `
		INSERT INTO setlist_songs(
			setlist_id,
			song_id,
			created_by,
			updated_by,
		) VALUES (
			$1, $2, $3, $4
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

	song := models.SetlistSong{}

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
		return models.SetlistSong{}, err
	}

	return song, nil
}

func SetlistSongsTableDeleteSong(songID string, setlistID string) error {
	log.Println("- SetlistSongsTableDeleteSong")

	query := `
		DELETE FROM setlist_songs
		WHERE song_id = $1
			AND setlist_id = $2
	`

	_, err := DB.Exec(query, songID, setlistID)
	if err != nil {
		return err
	}

	return nil
}

func SetlistsSongsTableGetAllSongsBySetlistID(setlistID string) ([]models.Setlist, error) {
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
		return []models.Setlist{}, err
	}

	defer rows.Close()

	setlists := []models.Setlist{}

	for rows.Next() {
		song := models.Setlist{}

		err := rows.Scan(
			&song.ID,
			&song.SetlistID,
			&song.BandID,
			&song.Name,
			&song.Notes,
			&song.ArtworkPath,
			&song.CreatedAt,
			&song.CreatedBy,
			&song.UpdatedAt,
			&song.UpdatedBy,
		)
		if err != nil {
			log.Printf("\n   Unable to get songs from setlist_songs with setlistID :", setlistID, err)
			return []models.Setlist{}, err
		}
		setlists = append(setlists, song)
	}

	return setlists, nil
}

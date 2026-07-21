package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// func CreateSetlistsTable(db *sql.DB) error {
// 	log.Println("- CreateSetlist")

// 	query := `
// 	CREATE TABLE IF NOT EXISTS setlists (
// 		id SERIAL PRIMARY KEY,
// 		setlist_id TEXT UNIQUE NOT NULL,
// 		band_id TEXT NOT NULL REFERENCES bands(band_id),
// 		name TEXT NOT NULL,
// 		notes TEXT,
// 		last_updated_by TEXT,
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 	)
// 	`

// 	_, err := db.Exec(query)
// 	if err != nil {
// 		log.Println("   Unable to create setlists table: ", err)
// 		log.Fatal(err)
// 	}
// 	return nil
// }

// func CreateSetlistSongsTable(db *sql.DB) error {
// 	log.Println("- CreateSetlistSongsTable")

// 	query := `
// 	CREATE TABLE IF NOT EXISTS setlist_songs (
// 		id SERIAL PRIMARY KEY,
// 		setlist_id TEXT NOT NULL
// 			REFERENCES setlists(setlist_id)
// 			ON DELETE CASCADE,
// 		song_id TEXT NOT NULL
// 			REFERENCES songs(song_id)
// 			ON DELETE CASCADE,
// 		position INT NOT NULL,
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

// 		UNIQUE(setlist_id, position)
// 	);
// 	`

// 	_, err := db.Exec(query)
// 	if err != nil {
// 		log.Println("   Unable to create setlist_songs table: ", err)
// 		log.Fatal(err)
// 	}
// 	return nil
// }

func SetlistsTableGetSetlistsByBandID(bandID string) (map[string]models.Setlist, error) {
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

	setlists := make(map[string]models.Setlist)

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
		setlists[setlist.SetlistID] = setlist
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

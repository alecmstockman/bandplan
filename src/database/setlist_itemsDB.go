package database

import (
	"bandplan/src/models"
	"fmt"
	"log"
)

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

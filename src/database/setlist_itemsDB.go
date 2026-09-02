package database

import (
	"bandplan/src/models"
	"errors"
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
		log.Println("   err deleting transition: ", err)
		return err
	}

	return tx.Commit()

}

func SetlistItemsTableDeleteBreak(breakID string, position int, setlistID string) error {
	log.Println("- SetlistItemsTableDeleteTransition")

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
		log.Println("   Error deleting break: ", err)
		return err
	}

	return tx.Commit()
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
	log.Println("- SetlistItemsGetItem")

	fmt.Println("itemType: ", itemType)
	fmt.Println("itemID: ", itemID)

	var query string

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
				setlist_id,
				item_type,
				break_id,
				position,
				pause_after_seconds,
				created_at, 
				created_by,
				updated_at,
				updated_by
			FROM setlist_items
			WHERE setlist_id = $1
				AND break_id = $2
			`
	}

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

func SetlistItemsGetSetlistOrder(setlistID string) ([]models.ReorderItem, error) {
	log.Println("- SetlistItemsGetSetlistOrder")

	query := `
		SELECT
			id,
			item_type,
			COALESCE(song_id, ''),
			COALESCE(transition_id, ''),
			COALESCE(break_id, ''),
			position
		FROM setlist_items
		WHERE setlist_id = $1
		ORDER BY position
	`

	rows, err := DB.Query(query, setlistID)
	if err != nil {
		log.Println("Unable to get setlist order from database: ", err)
		return []models.ReorderItem{}, err
	}

	defer rows.Close()

	var order []models.ReorderItem

	for rows.Next() {
		var orderItem models.ReorderItem

		var itemType string
		var songID string
		var transitionID string
		var breakID string

		err := rows.Scan(
			&orderItem.ItemID,
			&itemType,
			&songID,
			&transitionID,
			&breakID,
			&orderItem.Position,
		)
		if err != nil {
			log.Println("   Unable to get reorder item from setlist: ", err)
			return []models.ReorderItem{}, err
		}

		order = append(order, orderItem)
	}

	return order, nil
}

func SetlistItemsUpdateItemPosition(setlistID string, itemID string, newPosition int) error {
	log.Println("- SetlistItemsUpdateItemPosition")

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	maxQuery := `
		SELECT
			MAX(position)
		FROM setlist_items
		WHERE setlist_id = $1
	`

	var maxPositon int

	err = tx.QueryRow(maxQuery, setlistID).Scan(
		&maxPositon,
	)

	if newPosition >= maxPositon {
		return nil
	}

	return nil
}

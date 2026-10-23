package database

import (
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
)

type queryRower interface {
	QueryRow(query string, args ...any) *sql.Row
}

func setlistItemSource(itemType models.SetlistItemType) (string, string, error) {
	switch itemType {
	case models.SetlistItemSong:
		return "songs", "song_id", nil
	case models.SetlistItemTransition:
		return "transitions", "transition_id", nil
	case models.SetlistItemBreak:
		return "breaks", "break_id", nil
	default:
		return "", "", fmt.Errorf("invalid setlist item type: %s", itemType)
	}
}

func lockAuthorizedSetlist(tx *sql.Tx, setlistID, bandID, userID string) error {
	query := `
		SELECT s.setlist_id
		FROM setlists s
		JOIN band_members bm
			ON bm.band_id = s.band_id
			AND bm.user_id = $3
		WHERE s.setlist_id = $1
			AND s.band_id = $2
		FOR UPDATE OF s
	`

	var authorizedSetlistID string
	if err := tx.QueryRow(query, setlistID, bandID, userID).Scan(&authorizedSetlistID); err != nil {
		return fmt.Errorf("lockAuthorizedSetlist - authorize setlist:  %w", err)
	}

	return nil
}

func setlistItemsTableSaveItem(q queryRower, itemType models.SetlistItemType, itemID, userID, setlistID, bandID string) (models.SetlistItem, error) {
	itemTable, itemColumn, err := setlistItemSource(itemType)
	if err != nil {
		return models.SetlistItem{}, err
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
		SELECT
			$1,
			i.%s,
			COALESCE(
				(
					SELECT MAX(position) + 1
					FROM setlist_items
					WHERE setlist_id = $1
				),
				1
			),
			$3,
			$3,
			$4
		FROM %s i
		WHERE i.%s = $2
			AND i.band_id = $5
		RETURNING
			id,
			setlist_id,
			%s,
			position,
			created_at,
			created_by,
			updated_at,
			updated_by
	`, itemColumn, itemColumn, itemTable, itemColumn, itemColumn)

	item := models.SetlistItem{ItemType: itemType}
	err = q.QueryRow(
		query,
		setlistID,
		itemID,
		userID,
		string(itemType),
		bandID,
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
		return models.SetlistItem{}, fmt.Errorf("setlistItemsTableSaveItem - save setlist item:  %w", err)
	}

	return item, nil
}

func SetlistItemsTableSaveItem(itemType models.SetlistItemType, itemID, userID, setlistID, bandID string) (models.SetlistItem, error) {
	tx, err := DB.Begin()
	if err != nil {
		return models.SetlistItem{}, fmt.Errorf("SetlistItemsTableSaveItem - begin transaction:  %w", err)
	}
	defer tx.Rollback()

	if err := lockAuthorizedSetlist(tx, setlistID, bandID, userID); err != nil {
		return models.SetlistItem{}, err
	}

	item, err := setlistItemsTableSaveItem(tx, itemType, itemID, userID, setlistID, bandID)
	if err != nil {
		return models.SetlistItem{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.SetlistItem{}, fmt.Errorf("SetlistItemsTableSaveItem - commit transaction:  %w", err)
	}

	return item, nil
}

func setlistItemsTableDeleteItem(itemType models.SetlistItemType, itemID string, position int, setlistID, bandID, userID string) error {
	_, itemColumn, err := setlistItemSource(itemType)
	if err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("setlistItemsTableDeleteItem - begin transaction:  %w", err)
	}
	defer tx.Rollback()

	if err := lockAuthorizedSetlist(tx, setlistID, bandID, userID); err != nil {
		return err
	}

	deleteQuery := fmt.Sprintf(`
		DELETE FROM setlist_items
		WHERE %s = $1
			AND position = $2
			AND setlist_id = $3
		RETURNING position
	`, itemColumn)

	var deletedPosition int
	if err := tx.QueryRow(deleteQuery, itemID, position, setlistID).Scan(&deletedPosition); err != nil {
		return fmt.Errorf("setlistItemsTableDeleteItem - delete setlist item:  %w", err)
	}

	updateQuery := `
		UPDATE setlist_items
		SET
			position = position - 1,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $3
		WHERE setlist_id = $1
			AND position > $2
	`

	if _, err := tx.Exec(updateQuery, setlistID, deletedPosition, userID); err != nil {
		return fmt.Errorf("setlistItemsTableDeleteItem - compact setlist item positions:  %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("setlistItemsTableDeleteItem - commit transaction:  %w", err)
	}
	return nil
}

func SetlistItemsTableDeleteSong(songID string, position int, setlistID, bandID, userID string) error {
	return setlistItemsTableDeleteItem(models.SetlistItemSong, songID, position, setlistID, bandID, userID)
}

func SetlistItemsTableDeleteTransition(transitionID string, position int, setlistID, bandID, userID string) error {
	return setlistItemsTableDeleteItem(models.SetlistItemTransition, transitionID, position, setlistID, bandID, userID)
}

func SetlistItemsTableDeleteBreak(breakID string, position int, setlistID, bandID, userID string) error {
	return setlistItemsTableDeleteItem(models.SetlistItemBreak, breakID, position, setlistID, bandID, userID)
}

func SetlistItemsUpdateItem(setlistID string, itemType models.SetlistItemType, itemID string, pauseAfter int, bandID, userID string) error {
	itemTable, itemColumn, err := setlistItemSource(itemType)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE setlist_items si
		SET
			pause_after_seconds = $1,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $5
		FROM setlists s
		JOIN band_members bm
			ON bm.band_id = s.band_id
				AND bm.user_id = $5
		WHERE si.setlist_id = s.setlist_id
			AND si.setlist_id = $2
			AND si.%s = $3
			AND s.band_id = $4
			AND EXISTS (
				SELECT 1
				FROM %s i
				WHERE i.%s = $3
					AND i.band_id = $4
			)
	`, itemColumn, itemTable, itemColumn)

	result, err := DB.Exec(query, pauseAfter, setlistID, itemID, bandID, userID)
	if err != nil {
		return fmt.Errorf("SetlistItemsUpdateItem - update setlist item:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SetlistItemsUpdateItem - confirm setlist item update:  %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func SetlistItemsUpdateOrder(setlistID string, newOrder []models.ReorderItem, bandID, userID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("SetlistItemsUpdateOrder - begin transaction:  %w", err)
	}
	defer tx.Rollback()

	if err := lockAuthorizedSetlist(tx, setlistID, bandID, userID); err != nil {
		return err
	}

	var itemCount int
	if err := tx.QueryRow(
		"SELECT COUNT(*) FROM setlist_items WHERE setlist_id = $1",
		setlistID,
	).Scan(&itemCount); err != nil {
		return fmt.Errorf("SetlistItemsUpdateOrder - count setlist items:  %w", err)
	}
	if len(newOrder) != itemCount {
		return errors.New("setlist order does not contain every item")
	}

	itemIDs := make(map[string]struct{}, len(newOrder))
	positions := make(map[int]struct{}, len(newOrder))
	for _, item := range newOrder {
		if item.Position < 1 || item.Position > itemCount {
			return errors.New("setlist order contains an invalid position")
		}
		if _, exists := itemIDs[item.ItemID]; exists {
			return errors.New("setlist order contains a duplicate item")
		}
		if _, exists := positions[item.Position]; exists {
			return errors.New("setlist order contains a duplicate position")
		}
		itemIDs[item.ItemID] = struct{}{}
		positions[item.Position] = struct{}{}
	}

	updateQuery := `
		UPDATE setlist_items
		SET
			position = $1,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $4
		WHERE id = $2
			AND setlist_id = $3
	`

	for i, item := range newOrder {
		result, err := tx.Exec(updateQuery, -(i + 1), item.ItemID, setlistID, userID)
		if err != nil {
			return fmt.Errorf("SetlistItemsUpdateOrder - stage setlist item order:  %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("SetlistItemsUpdateOrder - confirm staged setlist item order:  %w", err)
		}
		if rowsAffected != 1 {
			return sql.ErrNoRows
		}
	}

	for _, item := range newOrder {
		result, err := tx.Exec(updateQuery, item.Position, item.ItemID, setlistID, userID)
		if err != nil {
			return fmt.Errorf("SetlistItemsUpdateOrder - update setlist item order:  %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("SetlistItemsUpdateOrder - confirm setlist item order update:  %w", err)
		}
		if rowsAffected != 1 {
			return sql.ErrNoRows
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("SetlistItemsUpdateOrder - commit transaction:  %w", err)
	}
	return nil
}

func SetlistItemsGetItem(setlistID string, itemType models.SetlistItemType, itemID, bandID, userID string) (models.SetlistItem, error) {
	itemTable, itemColumn, err := setlistItemSource(itemType)
	if err != nil {
		return models.SetlistItem{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			si.id,
			si.setlist_id,
			si.item_type,
			si.%s,
			si.position,
			si.pause_after_seconds,
			si.created_at,
			si.created_by,
			si.updated_at,
			si.updated_by
		FROM setlist_items si
		JOIN setlists s
			ON s.setlist_id = si.setlist_id
			AND s.band_id = $3
		JOIN band_members bm
			ON bm.band_id = s.band_id
			AND bm.user_id = $4
		JOIN %s i
			ON i.%s = si.%s
			AND i.band_id = s.band_id
		WHERE si.setlist_id = $1
			AND si.%s = $2
	`, itemColumn, itemTable, itemColumn, itemColumn, itemColumn)

	var item models.SetlistItem
	err = DB.QueryRow(query, setlistID, itemID, bandID, userID).Scan(
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
		return models.SetlistItem{}, fmt.Errorf("SetlistItemsGetItem - get setlist item:  %w", err)
	}

	return item, nil
}

func SetlistItemsGetSetlistOrder(setlistID, bandID, userID string) ([]models.ReorderItem, error) {
	query := `
		SELECT
			si.id,
			si.item_type,
			COALESCE(si.song_id, ''),
			COALESCE(si.transition_id, ''),
			COALESCE(si.break_id, ''),
			si.position
		FROM setlist_items si
		JOIN setlists s
			ON s.setlist_id = si.setlist_id
			AND s.band_id = $2
		JOIN band_members bm
			ON bm.band_id = s.band_id
			AND bm.user_id = $3
		WHERE si.setlist_id = $1
		ORDER BY si.position
	`

	rows, err := DB.Query(query, setlistID, bandID, userID)
	if err != nil {
		return nil, fmt.Errorf("SetlistItemsGetSetlistOrder - get setlist order:  %w", err)
	}
	defer rows.Close()

	order := []models.ReorderItem{}
	for rows.Next() {
		var orderItem models.ReorderItem
		var itemType, songID, transitionID, breakID string
		if err := rows.Scan(
			&orderItem.ItemID,
			&itemType,
			&songID,
			&transitionID,
			&breakID,
			&orderItem.Position,
		); err != nil {
			return nil, fmt.Errorf("SetlistItemsGetSetlistOrder - scan setlist order:  %w", err)
		}
		order = append(order, orderItem)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SetlistItemsGetSetlistOrder - iterate setlist order:  %w", err)
	}

	return order, nil
}

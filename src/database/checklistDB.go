package database

import (
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidChecklistReferences = errors.New("invalid checklist references")
var ErrChecklistNotFound = errors.New("checklist not found")
var ErrChecklistItemNotFound = errors.New("checklist item not found")
var ErrInvalidChecklistItemReferences = errors.New("invalid checklist item references")

func ChecklistListsTableCreateChecklist(checklist models.Checklist, requesterUserID, currentBandID string) error {
	query := `
		INSERT INTO checklist_lists (
			checklist_id,
			user_id,
			band_id,
			name,
			description,
			song_id,
			setlist_id,
			event_id,
			assigned_to,
			due_date,
			due_time,
			due_timezone,
			created_by,
			updated_by
		)
		SELECT
			$1,
			$2,
			$3,
			$4,
			NULLIF($5, ''),
			NULLIF($6::text, ''),
			NULLIF($7::text, ''),
			NULLIF($8::text, ''),
			NULLIF($9::text, ''),
			NULLIF($10::text, '')::date,
			NULLIF($11::text, '')::time,
			NULLIF($12::text, ''),
			$13,
			$13
		WHERE EXISTS (
			SELECT 1
			FROM band_members requester
			WHERE requester.band_id = $14
				AND requester.user_id = $13
		)
			AND (
				($2::text = $13 AND $3::text IS NULL)
				OR ($2::text IS NULL AND $3::text = $14)
			)
			AND (
				NULLIF($6::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM songs s
					WHERE s.song_id = $6
						AND s.band_id = $14
				)
			)
			AND (
				NULLIF($7::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM setlists sl
					WHERE sl.setlist_id = $7
						AND sl.band_id = $14
				)
			)
			AND (
				NULLIF($8::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM events e
					WHERE e.event_id = $8
						AND e.band_id = $14
				)
			)
			AND (
				NULLIF($9::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM band_members assignee
					WHERE assignee.user_id = $9
						AND assignee.band_id = $14
				)
			)
	`

	dueDate := ""
	if checklist.DueDate != nil {
		dueDate = checklist.DueDate.Format("2006-01-02")
	}
	dueTime := ""
	if checklist.DueTime != nil {
		dueTime = checklist.DueTime.Format("15:04:05")
	}

	result, err := DB.Exec(
		query,
		checklist.ChecklistID,
		checklist.UserID,
		checklist.BandID,
		checklist.Name,
		checklist.Description,
		checklist.SongID,
		checklist.SetlistID,
		checklist.EventID,
		checklist.AssignedTo,
		dueDate,
		dueTime,
		checklist.DueTimezone,
		requesterUserID,
		currentBandID,
	)
	if err != nil {
		return fmt.Errorf("ChecklistListsTableCreateChecklist - create checklist:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ChecklistListsTableCreateChecklist - confirm checklist creation:  %w", err)
	}
	if rowsAffected != 1 {
		return ErrInvalidChecklistReferences
	}

	return nil
}

func ChecklistListsTableGetChecklists(userID, bandID string) ([]models.Checklist, error) {
	query := `
		SELECT
			cl.id,
			cl.checklist_id,
			cl.user_id,
			cl.band_id,
			cl.name,
			COALESCE(cl.description, ''),
			cl.song_id,
			cl.setlist_id,
			cl.event_id,
			cl.assigned_to,
			cl.due_date,
			cl.due_time,
			cl.due_timezone,
			cl.is_complete,
			cl.completed_at,
			cl.completed_by,
			cl.created_at,
			cl.created_by,
			cl.updated_at,
			COALESCE(cl.updated_by, ''),
			COUNT(ci.item_id),
			COUNT(ci.item_id) FILTER (WHERE ci.is_complete = TRUE)
		FROM checklist_lists cl
		LEFT JOIN checklist_items ci
			ON ci.checklist_id = cl.checklist_id
		WHERE cl.user_id = $1
			OR (
				cl.band_id = $2
				AND EXISTS (
					SELECT 1
					FROM band_members bm
					WHERE bm.band_id = cl.band_id
						AND bm.user_id = $1
				)
			)
		GROUP BY cl.id
		ORDER BY LOWER(cl.name), cl.id
	`

	rows, err := DB.Query(query, userID, bandID)
	if err != nil {
		return nil, fmt.Errorf("ChecklistListsTableGetChecklists - query: %w", err)
	}
	defer rows.Close()

	checklists := make([]models.Checklist, 0)
	for rows.Next() {
		var checklist models.Checklist
		if err := rows.Scan(
			&checklist.ID,
			&checklist.ChecklistID,
			&checklist.UserID,
			&checklist.BandID,
			&checklist.Name,
			&checklist.Description,
			&checklist.SongID,
			&checklist.SetlistID,
			&checklist.EventID,
			&checklist.AssignedTo,
			&checklist.DueDate,
			&checklist.DueTime,
			&checklist.DueTimezone,
			&checklist.IsComplete,
			&checklist.CompletedAt,
			&checklist.CompletedBy,
			&checklist.CreatedAt,
			&checklist.CreatedBy,
			&checklist.UpdatedAt,
			&checklist.UpdatedBy,
			&checklist.ItemCount,
			&checklist.CompletedCount,
		); err != nil {
			return nil, fmt.Errorf("ChecklistListsTableGetChecklists - scan: %w", err)
		}
		checklists = append(checklists, checklist)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ChecklistListsTableGetChecklists - rows: %w", err)
	}

	return checklists, nil
}

func ChecklistListsTableGetChecklistByID(checklistID, userID, bandID string) (models.Checklist, error) {
	query := `
		SELECT
			cl.id,
			cl.checklist_id,
			cl.user_id,
			cl.band_id,
			cl.name,
			COALESCE(cl.description, ''),
			cl.song_id,
			cl.setlist_id,
			cl.event_id,
			cl.assigned_to,
			cl.due_date,
			cl.due_time,
			cl.due_timezone,
			cl.is_complete,
			cl.completed_at,
			cl.completed_by,
			cl.created_at,
			cl.created_by,
			cl.updated_at,
			COALESCE(cl.updated_by, '')
		FROM checklist_lists cl
		WHERE cl.checklist_id = $1
			AND (
				cl.user_id = $2
				OR (
					cl.band_id = $3
					AND EXISTS (
						SELECT 1
						FROM band_members bm
						WHERE bm.band_id = cl.band_id
							AND bm.user_id = $2
					)
				)
			)
	`

	var checklist models.Checklist
	err := DB.QueryRow(query, checklistID, userID, bandID).Scan(
		&checklist.ID,
		&checklist.ChecklistID,
		&checklist.UserID,
		&checklist.BandID,
		&checklist.Name,
		&checklist.Description,
		&checklist.SongID,
		&checklist.SetlistID,
		&checklist.EventID,
		&checklist.AssignedTo,
		&checklist.DueDate,
		&checklist.DueTime,
		&checklist.DueTimezone,
		&checklist.IsComplete,
		&checklist.CompletedAt,
		&checklist.CompletedBy,
		&checklist.CreatedAt,
		&checklist.CreatedBy,
		&checklist.UpdatedAt,
		&checklist.UpdatedBy,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Checklist{}, ErrChecklistNotFound
	}
	if err != nil {
		return models.Checklist{}, fmt.Errorf("ChecklistListsTableGetChecklistByID - scan: %w", err)
	}

	return checklist, nil
}

func ChecklistListsTableUpdateChecklist(checklist models.Checklist, requesterUserID, currentBandID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ChecklistListsTableUpdateChecklist - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var checklistID string
	err = tx.QueryRow(`
		SELECT checklist_id
		FROM checklist_lists
		WHERE checklist_id = $1
			AND (
				user_id = $2
				OR (
					band_id = $3
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = $3 AND user_id = $2
					)
				)
			)
		FOR UPDATE
	`, checklist.ChecklistID, requesterUserID, currentBandID).Scan(&checklistID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrChecklistNotFound
	}
	if err != nil {
		return fmt.Errorf("ChecklistListsTableUpdateChecklist - authorize: %w", err)
	}

	dueDate, dueTime := checklistDueValues(checklist.DueDate, checklist.DueTime)
	result, err := tx.Exec(`
		UPDATE checklist_lists
		SET user_id = $2,
			band_id = $3,
			name = $4,
			description = NULLIF($5, ''),
			song_id = NULLIF($6::text, ''),
			setlist_id = NULLIF($7::text, ''),
			event_id = NULLIF($8::text, ''),
			assigned_to = NULLIF($9::text, ''),
			due_date = NULLIF($10::text, '')::date,
			due_time = NULLIF($11::text, '')::time,
			due_timezone = NULLIF($12::text, ''),
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $13
		WHERE checklist_id = $1
			AND (
				($2::text = $13 AND $3::text IS NULL)
				OR ($2::text IS NULL AND $3::text = $14)
			)
			AND (NULLIF($6::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM songs WHERE song_id = $6 AND band_id = $14
			))
			AND (NULLIF($7::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM setlists WHERE setlist_id = $7 AND band_id = $14
			))
			AND (NULLIF($8::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM events WHERE event_id = $8 AND band_id = $14
			))
			AND (NULLIF($9::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM band_members WHERE user_id = $9 AND band_id = $14
			))
	`,
		checklistID,
		checklist.UserID,
		checklist.BandID,
		checklist.Name,
		checklist.Description,
		checklist.SongID,
		checklist.SetlistID,
		checklist.EventID,
		checklist.AssignedTo,
		dueDate,
		dueTime,
		checklist.DueTimezone,
		requesterUserID,
		currentBandID,
	)
	if err != nil {
		return fmt.Errorf("ChecklistListsTableUpdateChecklist - update: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ChecklistListsTableUpdateChecklist - rows affected: %w", err)
	}
	if affected != 1 {
		return ErrInvalidChecklistReferences
	}
	if checklist.UserID != nil {
		if _, err := tx.Exec(`
			UPDATE checklist_items
			SET assigned_to = NULL,
				updated_at = CURRENT_TIMESTAMP,
				updated_by = $2
			WHERE checklist_id = $1
				AND assigned_to IS NOT NULL
		`, checklistID, requesterUserID); err != nil {
			return fmt.Errorf("ChecklistListsTableUpdateChecklist - clear item assignments: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ChecklistListsTableUpdateChecklist - commit: %w", err)
	}
	return nil
}

func ChecklistListsTableDuplicateChecklist(checklistID, requesterUserID, currentBandID string) (string, error) {
	tx, err := DB.Begin()
	if err != nil {
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var source models.Checklist
	err = tx.QueryRow(`
		SELECT
			user_id, band_id, name, COALESCE(description, ''),
			assigned_to, due_date, due_time, due_timezone
		FROM checklist_lists
		WHERE checklist_id = $1
			AND (
				user_id = $2
				OR (
					band_id = $3
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = $3 AND user_id = $2
					)
				)
			)
		FOR SHARE
	`, checklistID, requesterUserID, currentBandID).Scan(
		&source.UserID,
		&source.BandID,
		&source.Name,
		&source.Description,
		&source.AssignedTo,
		&source.DueDate,
		&source.DueTime,
		&source.DueTimezone,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrChecklistNotFound
	}
	if err != nil {
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - load checklist: %w", err)
	}

	rows, err := tx.Query(`
		SELECT
			id, item_id, checklist_id, name, position, COALESCE(body, ''),
			song_id, setlist_id, event_id, assigned_to,
			due_date, due_time, due_timezone, is_complete,
			completed_at, completed_by, created_at, created_by,
			updated_at, COALESCE(updated_by, '')
		FROM checklist_items
		WHERE checklist_id = $1
		ORDER BY position
	`, checklistID)
	if err != nil {
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - load items: %w", err)
	}

	items := make([]models.ChecklistItem, 0)
	for rows.Next() {
		var item models.ChecklistItem
		if err := scanChecklistItem(rows, &item); err != nil {
			rows.Close()
			return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - item rows: %w", err)
	}
	if err := rows.Close(); err != nil {
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - close item rows: %w", err)
	}

	nameRunes := []rune(source.Name)
	if len(nameRunes) > 113 {
		nameRunes = nameRunes[:113]
	}
	copyName := string(nameRunes) + " (Copy)"
	newChecklistID := uuid.NewString()
	dueDate, dueTime := checklistDueValues(source.DueDate, source.DueTime)
	_, err = tx.Exec(`
		INSERT INTO checklist_lists (
			checklist_id, user_id, band_id, name, description,
			assigned_to, due_date, due_time, due_timezone,
			created_by, updated_by
		)
		VALUES (
			$1, $2, $3, $4, NULLIF($5, ''),
			NULLIF($6::text, ''), NULLIF($7::text, '')::date,
			NULLIF($8::text, '')::time, NULLIF($9::text, ''),
			$10, $10
		)
	`,
		newChecklistID,
		source.UserID,
		source.BandID,
		copyName,
		source.Description,
		source.AssignedTo,
		dueDate,
		dueTime,
		source.DueTimezone,
		requesterUserID,
	)
	if err != nil {
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - insert checklist: %w", err)
	}

	for _, item := range items {
		itemDueDate, itemDueTime := checklistDueValues(item.DueDate, item.DueTime)
		_, err = tx.Exec(`
			INSERT INTO checklist_items (
				item_id, checklist_id, name, position, body,
				song_id, setlist_id, event_id, assigned_to,
				due_date, due_time, due_timezone, created_by, updated_by
			)
			VALUES (
				$1, $2, $3, $4, NULLIF($5, ''),
				$6, $7, $8, $9,
				NULLIF($10::text, '')::date,
				NULLIF($11::text, '')::time,
				$12, $13, $13
			)
		`,
			uuid.NewString(),
			newChecklistID,
			item.Name,
			item.Position,
			item.Body,
			item.SongID,
			item.SetlistID,
			item.EventID,
			item.AssignedTo,
			itemDueDate,
			itemDueTime,
			item.DueTimezone,
			requesterUserID,
		)
		if err != nil {
			return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - insert item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("ChecklistListsTableDuplicateChecklist - commit: %w", err)
	}
	return newChecklistID, nil
}

func ChecklistListsTableDeleteChecklist(checklistID, requesterUserID, currentBandID string) error {
	result, err := DB.Exec(`
		DELETE FROM checklist_lists cl
		WHERE cl.checklist_id = $1
			AND (
				cl.user_id = $2
				OR (
					cl.band_id = $3
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = cl.band_id AND user_id = $2
					)
				)
			)
	`, checklistID, requesterUserID, currentBandID)
	if err != nil {
		return fmt.Errorf("ChecklistListsTableDeleteChecklist - delete: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ChecklistListsTableDeleteChecklist - rows affected: %w", err)
	}
	if affected != 1 {
		return ErrChecklistNotFound
	}
	return nil
}

func ChecklistItemsTableGetItems(checklistID string) ([]models.ChecklistItem, error) {
	rows, err := DB.Query(`
		SELECT
			id,
			item_id,
			checklist_id,
			name,
			position,
			COALESCE(body, ''),
			song_id,
			setlist_id,
			event_id,
			assigned_to,
			due_date,
			due_time,
			due_timezone,
			is_complete,
			completed_at,
			completed_by,
			created_at,
			created_by,
			updated_at,
			COALESCE(updated_by, '')
		FROM checklist_items
		WHERE checklist_id = $1
		ORDER BY is_complete ASC, position ASC
	`, checklistID)
	if err != nil {
		return nil, fmt.Errorf("ChecklistItemsTableGetItems - query: %w", err)
	}
	defer rows.Close()

	items := make([]models.ChecklistItem, 0)
	for rows.Next() {
		var item models.ChecklistItem
		if err := scanChecklistItem(rows, &item); err != nil {
			return nil, fmt.Errorf("ChecklistItemsTableGetItems - scan: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ChecklistItemsTableGetItems - rows: %w", err)
	}

	return items, nil
}

func ChecklistItemsTableCreateItem(item models.ChecklistItem, requesterUserID, currentBandID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableCreateItem - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var checklistID string
	err = tx.QueryRow(`
		SELECT checklist_id
		FROM checklist_lists
		WHERE checklist_id = $1
			AND (
				user_id = $2
				OR (
					band_id = $3
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = $3 AND user_id = $2
					)
				)
			)
		FOR UPDATE
	`, item.ChecklistID, requesterUserID, currentBandID).Scan(&checklistID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrInvalidChecklistItemReferences
	}
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableCreateItem - authorize: %w", err)
	}

	dueDate, dueTime := checklistDueValues(item.DueDate, item.DueTime)
	result, err := tx.Exec(`
		INSERT INTO checklist_items (
			item_id, checklist_id, name, position, body,
			song_id, setlist_id, event_id, assigned_to,
			due_date, due_time, due_timezone, created_by, updated_by
		)
		SELECT
			$1, $2, $3,
			COALESCE((SELECT MAX(position) + 1 FROM checklist_items WHERE checklist_id = $2), 0),
			NULLIF($4, ''),
			NULLIF($5::text, ''),
			NULLIF($6::text, ''),
			NULLIF($7::text, ''),
			NULLIF($8::text, ''),
			NULLIF($9::text, '')::date,
			NULLIF($10::text, '')::time,
			NULLIF($11::text, ''),
			$12,
			$12
		WHERE (NULLIF($5::text, '') IS NULL OR EXISTS (
			SELECT 1 FROM songs WHERE song_id = $5 AND band_id = $13
		))
		AND (NULLIF($6::text, '') IS NULL OR EXISTS (
			SELECT 1 FROM setlists WHERE setlist_id = $6 AND band_id = $13
		))
		AND (NULLIF($7::text, '') IS NULL OR EXISTS (
			SELECT 1 FROM events WHERE event_id = $7 AND band_id = $13
		))
		AND (NULLIF($8::text, '') IS NULL OR EXISTS (
			SELECT 1 FROM band_members WHERE user_id = $8 AND band_id = $13
		))
	`,
		uuid.NewString(),
		checklistID,
		item.Name,
		item.Body,
		item.SongID,
		item.SetlistID,
		item.EventID,
		item.AssignedTo,
		dueDate,
		dueTime,
		item.DueTimezone,
		requesterUserID,
		currentBandID,
	)
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableCreateItem - insert: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableCreateItem - rows affected: %w", err)
	}
	if affected != 1 {
		return ErrInvalidChecklistItemReferences
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ChecklistItemsTableCreateItem - commit: %w", err)
	}
	return nil
}

func ChecklistItemsTableGetItemByID(itemID, checklistID, userID, bandID string) (models.ChecklistItem, error) {
	row := DB.QueryRow(`
		SELECT
			ci.id,
			ci.item_id,
			ci.checklist_id,
			ci.name,
			ci.position,
			COALESCE(ci.body, ''),
			ci.song_id,
			ci.setlist_id,
			ci.event_id,
			ci.assigned_to,
			ci.due_date,
			ci.due_time,
			ci.due_timezone,
			ci.is_complete,
			ci.completed_at,
			ci.completed_by,
			ci.created_at,
			ci.created_by,
			ci.updated_at,
			COALESCE(ci.updated_by, '')
		FROM checklist_items ci
		JOIN checklist_lists cl ON cl.checklist_id = ci.checklist_id
		WHERE ci.item_id = $1
			AND ci.checklist_id = $2
			AND (
				cl.user_id = $3
				OR (
					cl.band_id = $4
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = cl.band_id AND user_id = $3
					)
				)
			)
	`, itemID, checklistID, userID, bandID)

	var item models.ChecklistItem
	err := scanChecklistItem(row, &item)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ChecklistItem{}, ErrChecklistItemNotFound
	}
	if err != nil {
		return models.ChecklistItem{}, fmt.Errorf("ChecklistItemsTableGetItemByID - scan: %w", err)
	}
	return item, nil
}

func ChecklistItemsTableSetCompletion(itemID, checklistID, userID, bandID string, complete bool) error {
	result, err := DB.Exec(`
		UPDATE checklist_items ci
		SET is_complete = $5,
			completed_at = CASE WHEN $5 THEN CURRENT_TIMESTAMP END,
			completed_by = CASE WHEN $5 THEN $3 END,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $3
		FROM checklist_lists cl
		WHERE ci.item_id = $1
			AND ci.checklist_id = $2
			AND cl.checklist_id = ci.checklist_id
			AND (
				cl.user_id = $3
				OR (
					cl.band_id = $4
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = cl.band_id AND user_id = $3
					)
				)
			)
	`, itemID, checklistID, userID, bandID, complete)
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableSetCompletion - update: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableSetCompletion - rows affected: %w", err)
	}
	if affected != 1 {
		return ErrChecklistItemNotFound
	}
	return nil
}

func ChecklistItemsTableResetChecklist(checklistID, userID, bandID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableResetChecklist - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var authorizedChecklistID string
	err = tx.QueryRow(`
		SELECT checklist_id
		FROM checklist_lists
		WHERE checklist_id = $1
			AND (
				user_id = $2
				OR (
					band_id = $3
					AND EXISTS (
						SELECT 1
						FROM band_members
						WHERE band_id = $3
							AND user_id = $2
					)
				)
			)
		FOR UPDATE
	`, checklistID, userID, bandID).Scan(&authorizedChecklistID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrChecklistNotFound
	}
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableResetChecklist - authorize: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE checklist_lists
		SET is_complete = FALSE,
			completed_at = NULL,
			completed_by = NULL,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $2
		WHERE checklist_id = $1
	`, authorizedChecklistID, userID)
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableResetChecklist - reset checklist: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE checklist_items
		SET is_complete = FALSE,
			completed_at = NULL,
			completed_by = NULL,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $2
		WHERE checklist_id = $1
			AND is_complete = TRUE
	`, authorizedChecklistID, userID)
	if err != nil {
		return fmt.Errorf("ChecklistItemsTableResetChecklist - reset items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ChecklistItemsTableResetChecklist - commit: %w", err)
	}
	return nil
}

func ChecklistItemsTableUpdateItem(item models.ChecklistItem, requesterUserID, currentBandID string, completion *bool) (models.ChecklistItem, error) {
	tx, err := DB.Begin()
	if err != nil {
		return models.ChecklistItem{}, fmt.Errorf("ChecklistItemsTableUpdateItem - begin transaction: %w", err)
	}
	defer tx.Rollback()

	var itemID string
	err = tx.QueryRow(`
		SELECT ci.item_id
		FROM checklist_items ci
		JOIN checklist_lists cl ON cl.checklist_id = ci.checklist_id
		WHERE ci.item_id = $1
			AND ci.checklist_id = $2
			AND (
				cl.user_id = $3
				OR (
					cl.band_id = $4
					AND EXISTS (
						SELECT 1 FROM band_members
						WHERE band_id = cl.band_id AND user_id = $3
					)
				)
			)
		FOR UPDATE OF ci
	`, item.ItemID, item.ChecklistID, requesterUserID, currentBandID).Scan(&itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ChecklistItem{}, ErrChecklistItemNotFound
	}
	if err != nil {
		return models.ChecklistItem{}, fmt.Errorf("ChecklistItemsTableUpdateItem - authorize: %w", err)
	}

	dueDate, dueTime := checklistDueValues(item.DueDate, item.DueTime)
	row := tx.QueryRow(`
		UPDATE checklist_items
		SET name = $3,
			body = NULLIF($4, ''),
			song_id = NULLIF($5::text, ''),
			setlist_id = NULLIF($6::text, ''),
			event_id = NULLIF($7::text, ''),
			assigned_to = NULLIF($8::text, ''),
			due_date = NULLIF($9::text, '')::date,
			due_time = NULLIF($10::text, '')::time,
			due_timezone = NULLIF($11::text, ''),
			is_complete = COALESCE($14::boolean, is_complete),
			completed_at = CASE
				WHEN $14 IS NULL OR $14 = is_complete THEN completed_at
				WHEN $14 THEN CURRENT_TIMESTAMP
			END,
			completed_by = CASE
				WHEN $14 IS NULL OR $14 = is_complete THEN completed_by
				WHEN $14 THEN $12
			END,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $12
		WHERE item_id = $1
			AND checklist_id = $2
			AND (NULLIF($5::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM songs WHERE song_id = $5 AND band_id = $13
			))
			AND (NULLIF($6::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM setlists WHERE setlist_id = $6 AND band_id = $13
			))
			AND (NULLIF($7::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM events WHERE event_id = $7 AND band_id = $13
			))
			AND (NULLIF($8::text, '') IS NULL OR EXISTS (
				SELECT 1 FROM band_members WHERE user_id = $8 AND band_id = $13
			))
		RETURNING
			id, item_id, checklist_id, name, position, COALESCE(body, ''),
			song_id, setlist_id, event_id, assigned_to,
			due_date, due_time, due_timezone, is_complete,
			completed_at, completed_by, created_at, created_by,
			updated_at, COALESCE(updated_by, '')
	`,
		itemID,
		item.ChecklistID,
		item.Name,
		item.Body,
		item.SongID,
		item.SetlistID,
		item.EventID,
		item.AssignedTo,
		dueDate,
		dueTime,
		item.DueTimezone,
		requesterUserID,
		currentBandID,
		completion,
	)

	var updated models.ChecklistItem
	err = scanChecklistItem(row, &updated)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ChecklistItem{}, ErrInvalidChecklistItemReferences
	}
	if err != nil {
		return models.ChecklistItem{}, fmt.Errorf("ChecklistItemsTableUpdateItem - update: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return models.ChecklistItem{}, fmt.Errorf("ChecklistItemsTableUpdateItem - commit: %w", err)
	}
	return updated, nil
}

type checklistItemScanner interface {
	Scan(dest ...any) error
}

func scanChecklistItem(scanner checklistItemScanner, item *models.ChecklistItem) error {
	return scanner.Scan(
		&item.ID,
		&item.ItemID,
		&item.ChecklistID,
		&item.Name,
		&item.Position,
		&item.Body,
		&item.SongID,
		&item.SetlistID,
		&item.EventID,
		&item.AssignedTo,
		&item.DueDate,
		&item.DueTime,
		&item.DueTimezone,
		&item.IsComplete,
		&item.CompletedAt,
		&item.CompletedBy,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.UpdatedAt,
		&item.UpdatedBy,
	)
}

func checklistDueValues(dueDate, dueTime *time.Time) (string, string) {
	dateValue := ""
	if dueDate != nil {
		dateValue = dueDate.Format("2006-01-02")
	}
	timeValue := ""
	if dueTime != nil {
		timeValue = dueTime.Format("15:04:05")
	}
	return dateValue, timeValue
}

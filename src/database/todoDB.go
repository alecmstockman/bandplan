package database

import (
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrInvalidTodoListReferences = errors.New("invalid todo list references")
var ErrInvalidTodoItemReferences = errors.New("invalid todo item references")

func TodoListsTableCreateTodoList(todoList models.ToDoList, requesterUserID, currentBandID string) error {
	query := `
		INSERT INTO todo_lists (
			todo_list_id,
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
			NULLIF($2, ''),
			NULLIF($3, ''),
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
				($2 = $13 AND $3 = '')
				OR ($2 = '' AND $3 = $14)
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
	if todoList.DueDate != nil {
		dueDate = todoList.DueDate.Format("2006-01-02")
	}
	dueTime := ""
	if todoList.DueTime != nil {
		dueTime = todoList.DueTime.Format("15:04:05")
	}

	result, err := DB.Exec(
		query,
		uuid.NewString(),
		todoList.UserID,
		todoList.BandID,
		todoList.Name,
		todoList.Description,
		todoList.SongID,
		todoList.SetlistID,
		todoList.EventID,
		todoList.AssignedTo,
		dueDate,
		dueTime,
		todoList.DueTimezone,
		requesterUserID,
		currentBandID,
	)
	if err != nil {
		return fmt.Errorf("create todo list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("confirm todo list creation: %w", err)
	}
	if rowsAffected != 1 {
		return ErrInvalidTodoListReferences
	}

	return nil
}

func TodoListsTableGetListsByUserID(userID string) ([]models.ToDoList, error) {
	fmt.Println("- TodoListsTableGetListsByUserID")

	query := `
			SELECT
				id,
				todo_list_id,
				COALESCE(user_id, ''),
				COALESCE(band_id, ''),
				name,
				description,
				is_primary,
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
				updated_by
			FROM todo_lists
			WHERE user_id = $1
		`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return []models.ToDoList{}, err
	}

	defer rows.Close()

	var lists []models.ToDoList

	for rows.Next() {
		var list models.ToDoList

		err := rows.Scan(
			&list.ID,
			&list.ToDoListID,
			&list.UserID,
			&list.BandID,
			&list.Name,
			&list.Description,
			&list.IsPrimaryUser,
			&list.SongID,
			&list.SetlistID,
			&list.EventID,
			&list.AssignedTo,
			&list.DueDate,
			&list.DueTime,
			&list.DueTimezone,
			&list.IsComplete,
			&list.CompletedAt,
			&list.CompletedBy,
			&list.CreatedAt,
			&list.CreatedBy,
			&list.UpdatedAt,
			&list.UpdatedBy,
		)
		if err != nil {
			return []models.ToDoList{}, err
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func TodoListsTableGetListByID(listID string) (models.ToDoList, error) {
	fmt.Println("- TodoListsTableGetListByID")

	query := `
			SELECT
				id,
				todo_list_id,
				COALESCE(user_id, ''),
				COALESCE(band_id, ''),
				name,
				description,
				is_primary,
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
				updated_by
			FROM todo_lists
			WHERE todo_list_id = $1
		`

	var todoList models.ToDoList

	err := DB.QueryRow(
		query,
		listID,
	).Scan(
		&todoList.ID,
		&todoList.ToDoListID,
		&todoList.UserID,
		&todoList.BandID,
		&todoList.Name,
		&todoList.Description,
		&todoList.IsPrimaryUser,
		&todoList.SongID,
		&todoList.SetlistID,
		&todoList.EventID,
		&todoList.AssignedTo,
		&todoList.DueDate,
		&todoList.DueTime,
		&todoList.DueTimezone,
		&todoList.IsComplete,
		&todoList.CompletedAt,
		&todoList.CompletedBy,
		&todoList.CreatedAt,
		&todoList.CreatedBy,
		&todoList.UpdatedAt,
		&todoList.UpdatedBy,
	)
	if err != nil {
		return models.ToDoList{}, err
	}

	return todoList, nil
}

func TodoListsTableGetPrimaryUserListByUserID(userID string) (models.ToDoList, error) {
	fmt.Println("- TodoListsTableGetPrimaryUserListByUserID")

	query := `
			SELECT
				id,
				todo_list_id,
				COALESCE(user_id, ''),
				COALESCE(band_id, ''),
				name,
				description,
				is_primary,
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
				updated_by
			FROM todo_lists
			WHERE user_id = $1
				AND is_primary = TRUE
		`

	var todoList models.ToDoList

	err := DB.QueryRow(
		query,
		userID,
	).Scan(
		&todoList.ID,
		&todoList.ToDoListID,
		&todoList.UserID,
		&todoList.BandID,
		&todoList.Name,
		&todoList.Description,
		&todoList.IsPrimaryUser,
		&todoList.SongID,
		&todoList.SetlistID,
		&todoList.EventID,
		&todoList.AssignedTo,
		&todoList.DueDate,
		&todoList.DueTime,
		&todoList.DueTimezone,
		&todoList.IsComplete,
		&todoList.CompletedAt,
		&todoList.CompletedBy,
		&todoList.CreatedAt,
		&todoList.CreatedBy,
		&todoList.UpdatedAt,
		&todoList.UpdatedBy,
	)
	if err != nil {
		return models.ToDoList{}, err
	}

	return todoList, nil
}

func TodoListsTableGetPrimaryBandListByBandID(bandID string) (models.ToDoList, error) {
	fmt.Println("- TodoListsTableGetPrimaryBandListByBandID")

	query := `
			SELECT
				id,
				todo_list_id,
				COALESCE(user_id, ''),
				COALESCE(band_id, ''),
				name,
				description,
				is_primary,
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
				updated_by
			FROM todo_lists t
			WHERE band_id = $1
				AND is_primary = TRUE
		`

	var todoList models.ToDoList

	err := DB.QueryRow(
		query,
		bandID,
	).Scan(
		&todoList.ID,
		&todoList.ToDoListID,
		&todoList.UserID,
		&todoList.BandID,
		&todoList.Name,
		&todoList.Description,
		&todoList.IsPrimaryUser,
		&todoList.SongID,
		&todoList.SetlistID,
		&todoList.EventID,
		&todoList.AssignedTo,
		&todoList.DueDate,
		&todoList.DueTime,
		&todoList.DueTimezone,
		&todoList.IsComplete,
		&todoList.CompletedAt,
		&todoList.CompletedBy,
		&todoList.CreatedAt,
		&todoList.CreatedBy,
		&todoList.UpdatedAt,
		&todoList.UpdatedBy,
	)
	if err != nil {
		return models.ToDoList{}, err
	}

	return todoList, nil
}

func TodoItemsTableGetItemsByListID(listID string) ([]models.ToDoItem, error) {
	query := `
		SELECT
			id,
			item_id,
			todo_list_id,
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
		FROM todo_items
		WHERE todo_list_id = $1
		ORDER BY position
	`

	rows, err := DB.Query(query, listID)
	if err != nil {
		return []models.ToDoItem{}, err
	}
	defer rows.Close()

	var items []models.ToDoItem

	for rows.Next() {
		var item models.ToDoItem

		err := rows.Scan(
			&item.ID,
			&item.ItemID,
			&item.ToDoListID,
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
		if err != nil {
			return []models.ToDoItem{}, err
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func TodoItemsTableCreateItem(item models.ToDoItem, requesterUserID, currentBandID string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("begin todo item creation: %w", err)
	}
	defer tx.Rollback()

	lockQuery := `
		SELECT todo_list_id
		FROM todo_lists
		WHERE todo_list_id = $1
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
	`

	var listID string
	err = tx.QueryRow(lockQuery, item.ToDoListID, requesterUserID, currentBandID).Scan(&listID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrInvalidTodoItemReferences
	}
	if err != nil {
		return fmt.Errorf("authorize todo item creation: %w", err)
	}

	query := `
		INSERT INTO todo_items (
			item_id,
			todo_list_id,
			name,
			position,
			body,
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
			COALESCE(
				(
					SELECT MAX(position) + 1
					FROM todo_items
					WHERE todo_list_id = $2
				),
				0
			),
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
		WHERE (
			NULLIF($5::text, '') IS NULL
			OR EXISTS (
				SELECT 1
				FROM songs
				WHERE song_id = $5
					AND band_id = $13
			)
		)
			AND (
				NULLIF($6::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM setlists
					WHERE setlist_id = $6
						AND band_id = $13
				)
			)
			AND (
				NULLIF($7::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM events
					WHERE event_id = $7
						AND band_id = $13
				)
			)
			AND (
				NULLIF($8::text, '') IS NULL
				OR EXISTS (
					SELECT 1
					FROM band_members
					WHERE user_id = $8
						AND band_id = $13
				)
			)
	`

	dueDate := ""
	if item.DueDate != nil {
		dueDate = item.DueDate.Format("2006-01-02")
	}
	dueTime := ""
	if item.DueTime != nil {
		dueTime = item.DueTime.Format("15:04:05")
	}

	result, err := tx.Exec(
		query,
		uuid.NewString(),
		listID,
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
		return fmt.Errorf("create todo item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("confirm todo item creation: %w", err)
	}
	if rowsAffected != 1 {
		return ErrInvalidTodoItemReferences
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit todo item creation: %w", err)
	}

	return nil
}

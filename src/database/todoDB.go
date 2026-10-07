package database

import (
	"bandplan/src/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrInvalidTodoListReferences = errors.New("invalid todo list references")

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

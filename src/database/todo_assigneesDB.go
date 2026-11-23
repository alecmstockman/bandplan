package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func ToDoListAssigneesTableCreateAssignees(ctx context.Context, todoListID, bandID, userID string, assigneeIDs []string) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ToDoListAssigneesTableCreateAssignees - begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := todoListAssigneesCreateInTransaction(ctx, tx, todoListID, bandID, userID, assigneeIDs); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ToDoListAssigneesTableCreateAssignees - commit: %w", err)
	}

	return nil
}

func todoListAssigneesCreateInTransaction(
	ctx context.Context,
	tx *sql.Tx,
	todoListID, bandID, userID string,
	assigneeIDs []string,
) error {
	query := `
		INSERT INTO todo_list_assignees(
			todo_list_assignee_id,
			user_id,
			band_id,
			todo_list_id,
			created_by,
			updated_by
		)
		SELECT $1, $2, $3, $4, $5, $5
		WHERE EXISTS (
			SELECT 1
			FROM band_members
			WHERE band_id = $3
				AND user_id = $2
		)
	`

	for _, assigneeID := range assigneeIDs {
		result, err := tx.ExecContext(
			ctx,
			query,
			uuid.NewString(),
			assigneeID,
			bandID,
			todoListID,
			userID,
		)
		if err != nil {
			return fmt.Errorf("ToDoListAssigneesTableCreateAssignees - insert assignee: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("ToDoListAssigneesTableCreateAssignees - rows affected: %w", err)
		}
		if rowsAffected != 1 {
			return fmt.Errorf("ToDoListAssigneesTableCreateAssignees - invalid assignee: %w", ErrInvalidTodoListReferences)
		}
	}

	return nil
}

func todoListAssigneesReplaceInTransaction(
	ctx context.Context,
	tx *sql.Tx,
	todoListID, bandID, userID string,
	assigneeIDs []string,
) error {
	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM todo_list_assignees WHERE todo_list_id = $1`,
		todoListID,
	); err != nil {
		return fmt.Errorf("replace to-do list assignees - delete existing: %w", err)
	}

	if len(assigneeIDs) == 0 {
		return nil
	}

	return todoListAssigneesCreateInTransaction(
		ctx,
		tx,
		todoListID,
		bandID,
		userID,
		assigneeIDs,
	)
}

func ToDoListAssigneesTableGetAssigneeIDs(
	ctx context.Context,
	todoListID, bandID string,
) (map[string]bool, error) {
	rows, err := DB.QueryContext(
		ctx,
		`
			SELECT user_id
			FROM todo_list_assignees
			WHERE todo_list_id = $1
				AND band_id = $2
		`,
		todoListID,
		bandID,
	)
	if err != nil {
		return nil, fmt.Errorf("get to-do list assignees: %w", err)
	}
	defer rows.Close()

	assigneeIDs := make(map[string]bool)
	for rows.Next() {
		var assigneeID string
		if err := rows.Scan(&assigneeID); err != nil {
			return nil, fmt.Errorf("scan to-do list assignee: %w", err)
		}
		assigneeIDs[assigneeID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate to-do list assignees: %w", err)
	}

	return assigneeIDs, nil
}

func todoItemAssigneesCreateInTransaction(
	ctx context.Context,
	tx *sql.Tx,
	itemID, bandID, userID string,
	assigneeIDs []string,
) error {
	query := `
		INSERT INTO todo_item_assignees(
			todo_item_assignee_id,
			user_id,
			band_id,
			item_id,
			created_by,
			updated_by
		)
		SELECT $1, $2, $3, $4, $5, $5
		WHERE EXISTS (
			SELECT 1
			FROM band_members
			WHERE band_id = $3
				AND user_id = $2
		)
	`

	for _, assigneeID := range assigneeIDs {
		result, err := tx.ExecContext(
			ctx,
			query,
			uuid.NewString(),
			assigneeID,
			bandID,
			itemID,
			userID,
		)
		if err != nil {
			return fmt.Errorf("create to-do item assignee: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("count created to-do item assignees: %w", err)
		}
		if rowsAffected != 1 {
			return fmt.Errorf("invalid to-do item assignee: %w", ErrInvalidTodoItemReferences)
		}
	}

	return nil
}

func todoItemAssigneesReplaceInTransaction(
	ctx context.Context,
	tx *sql.Tx,
	itemID, bandID, userID string,
	assigneeIDs []string,
) error {
	if _, err := tx.ExecContext(
		ctx,
		`DELETE FROM todo_item_assignees WHERE item_id = $1`,
		itemID,
	); err != nil {
		return fmt.Errorf("replace to-do item assignees - delete existing: %w", err)
	}

	if len(assigneeIDs) == 0 {
		return nil
	}

	return todoItemAssigneesCreateInTransaction(ctx, tx, itemID, bandID, userID, assigneeIDs)
}

func ToDoItemAssigneesTableGetAssigneeIDs(
	ctx context.Context,
	itemID, bandID string,
) (map[string]bool, error) {
	rows, err := DB.QueryContext(
		ctx,
		`
			SELECT user_id
			FROM todo_item_assignees
			WHERE item_id = $1
				AND band_id = $2
		`,
		itemID,
		bandID,
	)
	if err != nil {
		return nil, fmt.Errorf("get to-do item assignees: %w", err)
	}
	defer rows.Close()

	assigneeIDs := make(map[string]bool)
	for rows.Next() {
		var assigneeID string
		if err := rows.Scan(&assigneeID); err != nil {
			return nil, fmt.Errorf("scan to-do item assignee: %w", err)
		}
		assigneeIDs[assigneeID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate to-do item assignees: %w", err)
	}

	return assigneeIDs, nil
}

func ToDoItemAssigneesTableGetByListID(
	ctx context.Context,
	listID, bandID string,
) (map[string]map[string]bool, error) {
	rows, err := DB.QueryContext(
		ctx,
		`
			SELECT tia.item_id, tia.user_id
			FROM todo_item_assignees tia
			INNER JOIN todo_items ti
				ON ti.item_id = tia.item_id
			WHERE ti.todo_list_id = $1
				AND tia.band_id = $2
		`,
		listID,
		bandID,
	)
	if err != nil {
		return nil, fmt.Errorf("get to-do item assignees by list: %w", err)
	}
	defer rows.Close()

	assigneeIDsByItem := make(map[string]map[string]bool)
	for rows.Next() {
		var itemID, assigneeID string
		if err := rows.Scan(&itemID, &assigneeID); err != nil {
			return nil, fmt.Errorf("scan to-do item assignee by list: %w", err)
		}
		if assigneeIDsByItem[itemID] == nil {
			assigneeIDsByItem[itemID] = make(map[string]bool)
		}
		assigneeIDsByItem[itemID][assigneeID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate to-do item assignees by list: %w", err)
	}

	return assigneeIDsByItem, nil
}

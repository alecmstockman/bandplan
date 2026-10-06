package database

import (
	"bandplan/src/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrInvalidChecklistReferences = errors.New("invalid checklist references")

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
		uuid.NewString(),
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
		return fmt.Errorf("create checklist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("confirm checklist creation: %w", err)
	}
	if rowsAffected != 1 {
		return ErrInvalidChecklistReferences
	}

	return nil
}

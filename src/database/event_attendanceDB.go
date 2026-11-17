package database

import (
	"bandplan/src/models"
	"context"
	"fmt"
)

func EventsTableCreateEventWithAttendance(event models.Event, attendance []models.EventAttendance) (models.Event, error) {
	tx, err := DB.Begin()
	if err != nil {
		return models.Event{}, fmt.Errorf("EventsTableCreateEventWithAttendance - begin: %w", err)
	}
	defer tx.Rollback()

	event, err = eventsTableCreateEvent(tx, event)
	if err != nil {
		return models.Event{}, err
	}

	query := `
		INSERT INTO event_attendance (
			event_attendance_id,
			event_id,
			user_id,
			status
		) VALUES ($1, $2, $3, $4)
	`

	for _, attendee := range attendance {
		if _, err := tx.Exec(
			query,
			attendee.EventAttendanceID,
			event.EventID,
			attendee.UserID,
			attendee.Status,
		); err != nil {
			return models.Event{}, fmt.Errorf("EventsTableCreateEventWithAttendance - insert attendance: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Event{}, fmt.Errorf("EventsTableCreateEventWithAttendance - commit: %w", err)
	}

	return event, nil
}

func EventAttendanceGetByEventID(ctx context.Context, eventID, bandID string) ([]models.EventAttendance, error) {
	query := `
		SELECT
			ea.id,
			ea.event_attendance_id,
			ea.event_id,
			ea.user_id,
			u.display_name,
			ea.status,
			ea.created_at,
			ea.updated_at
		FROM event_attendance ea
		JOIN events e
			ON e.event_id = ea.event_id
		JOIN users u
			ON u.user_id = ea.user_id
		WHERE ea.event_id = $1
			AND e.band_id = $2
		ORDER BY
			CASE ea.status
				WHEN 'confirmed' THEN 1
				WHEN 'maybe' THEN 2
				WHEN 'cannot' THEN 3
				ELSE 4
			END,
			u.display_name
	`

	rows, err := DB.QueryContext(ctx, query, eventID, bandID)
	if err != nil {
		return nil, fmt.Errorf("EventAttendanceGetByEventID - query: %w", err)
	}
	defer rows.Close()

	attendance := []models.EventAttendance{}
	for rows.Next() {
		var attendee models.EventAttendance
		if err := rows.Scan(
			&attendee.ID,
			&attendee.EventAttendanceID,
			&attendee.EventID,
			&attendee.UserID,
			&attendee.DisplayName,
			&attendee.Status,
			&attendee.CreatedAt,
			&attendee.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("EventAttendanceGetByEventID - scan: %w", err)
		}
		attendance = append(attendance, attendee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("EventAttendanceGetByEventID - rows: %w", err)
	}

	return attendance, nil
}

func EventAttendanceUpdateStatus(ctx context.Context, eventID, userID, bandID string, status models.EventAttendanceStatus) (bool, error) {
	query := `
		UPDATE event_attendance ea
		SET
			status = $4,
			updated_at = CURRENT_TIMESTAMP
		FROM events e
		WHERE ea.event_id = e.event_id
			AND ea.event_id = $1
			AND ea.user_id = $2
			AND e.band_id = $3
	`

	result, err := DB.ExecContext(ctx, query, eventID, userID, bandID, status)
	if err != nil {
		return false, fmt.Errorf("EventAttendanceUpdateStatus - update: %w", err)
	}

	updated, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("EventAttendanceUpdateStatus - rows affected: %w", err)
	}

	return updated == 1, nil
}

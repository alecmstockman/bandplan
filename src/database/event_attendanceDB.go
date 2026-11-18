package database

import (
	"bandplan/src/models"
	"context"
	"fmt"

	"github.com/lib/pq"
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

func EventsTableUpdateEventWithAttendance(ctx context.Context, event models.Event, inviteeIDs []string) (models.Event, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Event{}, fmt.Errorf("EventsTableUpdateEventWithAttendance - begin: %w", err)
	}
	defer tx.Rollback()

	event, err = eventsTableUpdateEvent(tx, event)
	if err != nil {
		return models.Event{}, err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM event_attendance ea
		USING band_members bm
		WHERE ea.event_id = $1
			AND bm.band_id = $2
			AND bm.user_id = ea.user_id
			AND ea.user_id <> ALL($3::text[])
	`, event.EventID, event.BandID, pq.Array(inviteeIDs))
	if err != nil {
		return models.Event{}, fmt.Errorf("EventsTableUpdateEventWithAttendance - remove invitees: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO event_attendance (event_attendance_id, event_id, user_id, status)
		SELECT gen_random_uuid()::text, $1, bm.user_id, $4
		FROM band_members bm
		WHERE bm.band_id = $2
			AND bm.user_id = ANY($3::text[])
		ON CONFLICT (event_id, user_id) DO NOTHING
	`, event.EventID, event.BandID, pq.Array(inviteeIDs), models.EventAttendanceUnseen)
	if err != nil {
		return models.Event{}, fmt.Errorf("EventsTableUpdateEventWithAttendance - add invitees: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Event{}, fmt.Errorf("EventsTableUpdateEventWithAttendance - commit: %w", err)
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

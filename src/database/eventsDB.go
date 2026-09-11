package database

import (
	"bandplan/src/models"
	"log"
	"time"
)

func EventsTableCreateEvent(event models.Event) (models.Event, error) {
	log.Println("- EventsTableCreateEvent")

	query := `
		INSERT INTO events (
			event_id,
			band_id,
			name,
			slug,
			image_id,
			image_path,
			event_date,
			event_type,
			recurrence,
			location,
			address,
			start_time,
			end_time,
			time_zone,
			set_location,
			load_in_time,
			load_in_instructions,
			set_time,
			set_length_seconds,
			venue_name,
			address_one,
			address_two,
			city,
			state,
			zip_code,
			presale_ticket_price,
			ticket_price,
			ticket_link,
			setlist_id,
			notes,
			link_one_name,
			link_one,
			link_two_name,
			link_two,
			created_by,
			updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, NULLIF($26, '')::NUMERIC,
			NULLIF($27, '')::NUMERIC, $28, NULLIF($29, ''), $30,
			$31, $32, $33, $34, $35, $36
		)
		RETURNING id, created_at, updated_at
	`

	nullableTime := func(value *time.Time) any {
		if value.IsZero() {
			return nil
		}
		return value
	}

	err := DB.QueryRow(
		query,
		event.EventID,
		event.BandID,
		event.Name,
		event.Slug,
		event.ImageID,
		event.ImagePath,
		event.EventDate,
		event.EventType,
		event.Recurrence,
		event.Location,
		event.Address,
		nullableTime(event.StartTime),
		nullableTime(event.EndTime),
		event.Timezone,
		event.SetLocation,
		nullableTime(event.LoadInTime),
		event.LoadInInstructions,
		nullableTime(event.SetTime),
		event.SetLengthSeconds,
		event.VenueName,
		event.AddressOne,
		event.AddressTwo,
		event.City,
		event.State,
		event.ZipCode,
		event.PresaleTicketPrice,
		event.TicketPrice,
		event.TicketLink,
		event.SetlistID,
		event.Notes,
		event.LinkOneName,
		event.LinkOne,
		event.LinkTwoName,
		event.LinkTwo,
		event.CreatedBy,
		event.UpdatedBy,
	).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		log.Println("   Unable to create event: ", err)
		return models.Event{}, err
	}

	return event, nil
}

func EventsTableGetAllEventsByBandIDAndUserID(bandID string, userID string) ([]models.Event, error) {
	log.Println("- EventsTableGetAllEventsByBandID")

	query := `
		SELECT 
			e.id,
			e.event_id,
			e.band_id,
			e.name,
			e.slug,
			e.image_id,
			e.image_path,
			e.event_date,
			e.event_type,
			e.recurrence,
			COALESCE(e.location, ''),
			COALESCE(e.address, ''),
			e.start_time,
			e.end_time,
			e.time_zone,
			e.set_location,
			e.load_in_time,
			COALESCE(e.load_in_instructions, ''),
			e.set_time,
			COALESCE(e.set_length_seconds, 0),
			COALESCE(e.venue_name, ''),
			COALESCE(e.address_one, ''),
			COALESCE(e.address_two, ''),
			COALESCE(e.city, ''),
			COALESCE(e.state, ''),
			COALESCE(e.zip_code, ''),
			COALESCE(e.presale_ticket_price, 0),
			COALESCE(e.ticket_price, 0),
			COALESCE(e.ticket_link, ''),
			COALESCE(e.setlist_id, ''),
			COALESCE(e.notes, ''),
			COALESCE(e.link_one_name, ''),
			COALESCE(e.link_one, ''),
			COALESCE(e.link_two_name, ''), 
			COALESCE(e.link_two, ''),
			e.created_at,
			e.created_by,
			e.updated_at,
			COALESCE(e.updated_by, '')
		FROM events e
		WHERE band_id = $1
		AND EXISTS (
			SELECT 1 
			FROM band_members bm
			WHERE bm.band_id = e.band_id
			AND bm.user_id = $2
		)
	`

	rows, err := DB.Query(query, bandID, userID)
	if err != nil {
		log.Println("   Unable to get events by bandID and userID from database: ", err)
		return []models.Event{}, err
	}

	defer rows.Close()

	events := []models.Event{}

	for rows.Next() {

		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.EventID,
			&event.BandID,
			&event.Name,
			&event.Slug,
			&event.ImageID,
			&event.ImagePath,
			&event.EventDate,
			&event.EventType,
			&event.Recurrence,
			&event.Location,
			&event.Address,
			&event.StartTime,
			&event.EndTime,
			&event.Timezone,
			&event.SetLocation,
			&event.LoadInTime,
			&event.LoadInInstructions,
			&event.SetTime,
			&event.SetLengthSeconds,
			&event.VenueName,
			&event.AddressOne,
			&event.AddressTwo,
			&event.City,
			&event.State,
			&event.ZipCode,
			&event.PresaleTicketPrice,
			&event.TicketPrice,
			&event.TicketLink,
			&event.SetlistID,
			&event.Notes,
			&event.LinkOneName,
			&event.LinkOne,
			&event.LinkTwoName,
			&event.LinkTwo,
			&event.CreatedAt,
			&event.CreatedBy,
			&event.UpdatedAt,
			&event.UpdatedBy,
		)
		if err != nil {
			log.Println("   Unable to get event: ", err)
			return []models.Event{}, err
		}

		events = append(events, event)
	}

	return events, nil
}

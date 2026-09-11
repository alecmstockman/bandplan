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

	nullableTime := func(value time.Time) any {
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

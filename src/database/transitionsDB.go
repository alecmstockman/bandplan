package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func TransitionsTableCreateTransition(transition models.Transition) (models.Transition, error) {
	log.Println("- TransitionsTableCreateTransition")

	fmt.Printf("%+v\n", transition)

	transitionID := uuid.New().String()

	var newTransition models.Transition

	query := `
	INSERT INTO transitions(
		transition_id,
		band_id,
		title,
		title_slug,
		length_seconds,

		bpm,
		time_signature,
		musical_key,
		tuning,
		capo,

		explicit,
		chords,
		chart_link,
		link_one,
		link_two,

		link_three,
		lyrics,
		notes,
		created_by,
		updated_by
	)
	VALUES (
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20
	)
	RETURNING
		id, 
		transition_id,
		band_id,
		title,
		title_slug,
		length_seconds,
		bpm,
		time_signature,
		musical_key,
		tuning,
		capo,
		explicit,
		chords,
		chart_link,
		link_one,
		link_two,
		link_three,
		lyrics,
		notes,
		created_at,
		created_by,
		updated_at,
		updated_by
	`

	err := DB.QueryRow(
		query,
		transitionID,
		transition.BandID,
		transition.Title,
		transition.TitleSlug,
		transition.LengthSeconds,

		transition.BPM,
		transition.TimeSignature,
		transition.Key,
		transition.Tuning,
		transition.Capo,

		transition.Explicit,
		transition.Chords,
		transition.ChartLink,
		transition.LinkOne,
		transition.LinkTwo,

		transition.LinkThree,
		transition.Lyrics,
		transition.Notes,
		transition.CreatedBy,
		transition.UpdatedBy,
	).Scan(
		&newTransition.ID,
		&newTransition.TransitionID,
		&newTransition.BandID,
		&newTransition.Title,
		&newTransition.TitleSlug,
		&newTransition.LengthSeconds,
		&newTransition.BPM,
		&newTransition.TimeSignature,
		&newTransition.Key,
		&newTransition.Tuning,
		&newTransition.Capo,
		&newTransition.Explicit,
		&newTransition.Chords,
		&newTransition.ChartLink,
		&newTransition.LinkOne,
		&newTransition.LinkTwo,
		&newTransition.LinkThree,
		&newTransition.Lyrics,
		&newTransition.Notes,
		&newTransition.CreatedAt,
		&newTransition.CreatedBy,
		&newTransition.UpdatedAt,
		&newTransition.UpdatedBy,
	)

	if err != nil {
		log.Println("   Unable to save transition to db: ", err)
		return models.Transition{}, err
	}

	return newTransition, nil
}

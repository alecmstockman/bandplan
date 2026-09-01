package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func TransitionsTableCreateTransition(transition models.Transition) (models.Transition, error) {
	log.Println("- TransitionsTableCreateTransition")

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
		transition.Slug,
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
		&newTransition.Slug,
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

func TransitionsTableDeleteTransition(transitionID string) error {
	log.Println("- TransitionsTableDeleteTransition")

	query := `
	DELETE FROM transitions
	WHERE transition_id = $1
	`
	_, err := DB.Exec(query, transitionID)
	if err != nil {
		log.Printf("   Unable to delete transition; %v from songs db: %v\n", transitionID, err)
		return err
	}
	return nil
}

func TransitionsTableGetTransitionByID(transitionID string, bandID string) (models.Transition, error) {
	log.Println("- TransitionsTableGetTransitionByID")

	query := `
	SELECT
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

		COALESCE(link_one, ''),
		COALESCE(link_two, ''),
		COALESCE(link_three, ''),

		lyrics,
		notes,

		created_at,
		created_by,
		updated_at,
		updated_by
	FROM transitions
	WHERE transition_id = $1
		AND band_id = $2
	`

	var transition models.Transition

	err := DB.QueryRow(query, transitionID, bandID).Scan(
		&transition.ID,
		&transition.TransitionID,
		&transition.BandID,

		&transition.Title,
		&transition.Slug,
		&transition.LengthSeconds,
		&transition.BPM,
		&transition.TimeSignature,
		&transition.Key,
		&transition.Tuning,
		&transition.Capo,

		&transition.Explicit,
		&transition.Chords,
		&transition.ChartLink,

		&transition.LinkOne,
		&transition.LinkTwo,
		&transition.LinkThree,

		&transition.Lyrics,
		&transition.Notes,

		&transition.CreatedAt,
		&transition.CreatedBy,
		&transition.UpdatedAt,
		&transition.UpdatedBy,
	)
	if err != nil {
		log.Println("   Unable to get transition from transitions db: ", err)
		return models.Transition{}, err
	}

	return transition, nil
}

func TransitionsTableUpdateTransition(transition models.Transition) (bool, error) {
	log.Println("- TransitionsTableUpdateTransition")

	query := `
		UPDATE transitions
		SET
			title = $1,
			title_slug = $2,
			length_seconds = $3,
			bpm = $4,
			time_signature = $5,
			musical_key = $6,
			tuning = $7,
			capo = $8,
			explicit = $9,
			link_one = $10,
			lyrics = $11,
			notes = $12,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $13
		WHERE transition_id = $14
			AND band_id = $15
	`

	result, err := DB.Exec(
		query,
		transition.Title,
		transition.Slug,
		transition.LengthSeconds,
		transition.BPM,
		transition.TimeSignature,
		transition.Key,
		transition.Tuning,
		transition.Capo,
		transition.Explicit,
		transition.LinkOne,
		transition.Lyrics,
		transition.Notes,
		transition.UpdatedBy,
		transition.TransitionID,
		transition.BandID,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update transition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("unable to confirm transition update: %w", err)
	}

	return rowsAffected == 1, nil
}

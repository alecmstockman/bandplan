package database

import (
	"bandplan/src/models"
	"fmt"

	"github.com/google/uuid"
)

func transitionsTableCreateTransition(q queryRower, transition models.Transition, bandID, userID string) (models.Transition, error) {

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
	SELECT
		$1, $2, $3, $4, $5,
		$6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20
	FROM band_members bm
	WHERE bm.band_id = $2
		AND bm.user_id = $21
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

	err := q.QueryRow(
		query,
		transitionID,
		bandID,
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
		userID,
		userID,
		userID,
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
		return models.Transition{}, fmt.Errorf("transitionsTableCreateTransition - create transition:  %w", err)
	}

	return newTransition, nil
}

func TransitionsTableCreateTransition(transition models.Transition, bandID, userID string) (models.Transition, error) {
	return transitionsTableCreateTransition(DB, transition, bandID, userID)
}

func TransitionsTableCreateAndAttachTransition(transition models.Transition, setlistID, bandID, userID string) (models.Transition, error) {
	tx, err := DB.Begin()
	if err != nil {
		return models.Transition{}, fmt.Errorf("TransitionsTableCreateAndAttachTransition - begin transaction:  %w", err)
	}
	defer tx.Rollback()

	if err := lockAuthorizedSetlist(tx, setlistID, bandID, userID); err != nil {
		return models.Transition{}, err
	}

	newTransition, err := transitionsTableCreateTransition(tx, transition, bandID, userID)
	if err != nil {
		return models.Transition{}, err
	}

	if _, err := setlistItemsTableSaveItem(
		tx,
		models.SetlistItemTransition,
		newTransition.TransitionID,
		userID,
		setlistID,
		bandID,
	); err != nil {
		return models.Transition{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Transition{}, fmt.Errorf("TransitionsTableCreateAndAttachTransition - commit transaction:  %w", err)
	}

	return newTransition, nil
}

func TransitionsTableDeleteTransition(transitionID, bandID, userID string) (bool, error) {

	query := `
	DELETE FROM transitions
	WHERE transition_id = $1
		AND band_id = $2
		AND EXISTS (
			SELECT 1
			FROM band_members bm
			WHERE bm.band_id = transitions.band_id
				AND bm.user_id = $3
		)
	`
	result, err := DB.Exec(query, transitionID, bandID, userID)
	if err != nil {
		return false, fmt.Errorf("TransitionsTableDeleteTransition - delete transition:  %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("TransitionsTableDeleteTransition - confirm transition deletion:  %w", err)
	}
	return rowsAffected == 1, nil
}

func TransitionsTableGetTransitionByID(transitionID, bandID, userID string) (models.Transition, error) {

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
		AND EXISTS (
			SELECT 1
			FROM band_members bm
			WHERE bm.band_id = transitions.band_id
				AND bm.user_id = $3
		)
	`

	var transition models.Transition

	err := DB.QueryRow(query, transitionID, bandID, userID).Scan(
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
		return models.Transition{}, fmt.Errorf("TransitionsTableGetTransitionByID - get transition:  %w", err)
	}

	return transition, nil
}

func TransitionsTableUpdateTransition(transition models.Transition, bandID, userID string) (bool, error) {

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
			AND EXISTS (
				SELECT 1
				FROM band_members bm
				WHERE bm.band_id = transitions.band_id
					AND bm.user_id = $16
			)
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
		userID,
		transition.TransitionID,
		bandID,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("TransitionsTableUpdateTransition - update transition:  %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("TransitionsTableUpdateTransition - confirm transition update:  %w", err)
	}

	return rowsAffected == 1, nil
}

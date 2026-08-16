package database

import (
	"bandplan/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func BreaksTableCreateBreak(breakItem models.Break) (models.Break, error) {
	log.Println("- BreaksTableCreateBreak")

	breakID := uuid.New().String()

	var newBreak models.Break

	query := `
		INSERT INTO breaks (
			break_id,
			band_id,
			title,
			title_slug,
			notes,
			length_seconds,
			link_one,
			link_two,
			created_by,
			updated_by
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10
		)
		RETURNING
			id,
			break_id,
			band_id,
			title,
			title_slug,
			notes,
			length_seconds,
			link_one,
			link_two,
			created_at,
			created_by,
			updated_at,
			updated_by
	`

	err := DB.QueryRow(
		query,
		breakID,
		breakItem.BandID,
		breakItem.Title,
		breakItem.Slug,
		breakItem.Notes,
		breakItem.LengthSeconds,
		breakItem.LinkOne,
		breakItem.LinkTwo,
		breakItem.CreatedBy,
		breakItem.UpdatedBy,
	).Scan(
		&newBreak.ID,
		&newBreak.BreakID,
		&newBreak.BandID,
		&newBreak.Title,
		&newBreak.Slug,
		&newBreak.Notes,
		&newBreak.LengthSeconds,
		&newBreak.LinkOne,
		&newBreak.LinkTwo,
		&newBreak.CreatedAt,
		&newBreak.CreatedBy,
		&newBreak.UpdatedAt,
		&newBreak.UpdatedBy,
	)
	if err != nil {
		return models.Break{}, fmt.Errorf("unable to create break: %w", err)
	}

	return newBreak, nil
}

func BreaksTableDeleteBreak(breakID string) error {
	log.Println("- BreaksTableDeleteBreak")

	query := `
	DELETE FROM breaks
	WHERE break_id = $1
	`

	_, err := DB.Exec(query, breakID)
	if err != nil {
		log.Println("   Unable to delete break: ", err)
		return err
	}
	return nil
}

func BreaksTableGetBreakByID(breakID string) (models.Break, error) {
	log.Println("- BreaksTableGetBreakByID")

	query := `
	SELECT 
		id,
		break_id,
		band_id,
		title,
		title_slug,
		notes,
		length_seconds,
		link_one,
		link_two,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM breaks
	WHERE break_id = $1
	`

	var breakItem models.Break

	err := DB.QueryRow(query, breakID).Scan(
		&breakItem.ID,
		&breakItem.BreakID,
		&breakItem.BandID,
		&breakItem.Title,
		&breakItem.Slug,
		&breakItem.Notes,
		&breakItem.LengthSeconds,
		&breakItem.LinkOne,
		&breakItem.LinkTwo,
		&breakItem.CreatedAt,
		&breakItem.CreatedBy,
		&breakItem.UpdatedAt,
		&breakItem.UpdatedBy,
	)
	if err != nil {
		log.Println("   Unable to get break from database: ", err)
		return models.Break{}, err
	}

	return breakItem, nil
}

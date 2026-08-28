package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

func (s Service) SetlistCreate(ctx context.Context, user models.User, band models.Band, input CreateSetlistInput) (models.Setlist, error) {
	log.Println("- CreateSetlist")

	title := strings.TrimSpace(input.Title)
	notes := strings.TrimSpace(input.Notes)
	artworkPath := strings.TrimSpace(input.TempArtID)
	slug := helpers.MakeSlug(title)

	if input.TempArtID != "" {
		path, err := s.ServiceCreatePermSetlistImage(
			ctx,
			input.TempArtID,
			band.Slug,
			slug,
		)
		if err != nil {
			log.Println("   Unable to save temporary artwork image in R2: ", err)
		} else {
			artworkPath = path
		}
	}

	setlist := models.Setlist{
		BandID:      band.BandID,
		Name:        title,
		Slug:        slug,
		Explicit:    false,
		Notes:       notes,
		ArtworkID:   input.TempArtID,
		ArtworkPath: artworkPath,
		CreatedBy:   user.UserID,
		UpdatedBy:   user.UserID,
	}

	if err := database.SetlistsTableCreateSetlist(setlist); err != nil {
		log.Printf("   Unable to create setlist %v in setlists table: %v\n", title, err)
		return models.Setlist{}, err
	}

	return setlist, nil
}

func (s Service) DeleteTempImage(ctx context.Context, imageID string, bandSlug string) error {
	log.Println("- HelperDeleteTempImage")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID empty")
	}

	key := fmt.Sprintf(
		"temp-images/%s/%s.webp",
		bandSlug,
		imageID,
	)

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete temporary image from R2:", err)
		return errors.New("Unable to delete temporary image from R2")
	}

	return nil
}

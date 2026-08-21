package services

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"bandplan/src/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type SetlistService struct {
	DB      *sql.DB
	Storage *storage.R2Storage
}

type CreateSetlistInput struct {
	Title     string
	Notes     string
	TempArtID string
}

func (s SetlistService) SetlistCreate(ctx context.Context, user models.User, band models.Band, input CreateSetlistInput) (models.Setlist, error) {
	log.Println("- CreateSetlist")

	title := strings.TrimSpace(input.Title)
	notes := strings.TrimSpace(input.Notes)
	artworkPath := strings.TrimSpace(input.TempArtID)
	slug := helpers.MakeSlug(title)

	if input.TempArtID != "" {
		path, err := s.HelperCreatePermSetlistImage(
			ctx,
			input.TempArtID,
			band.Slug,
			slug,
		)
		if err != nil {
			// decide whether this should fail creation or just continue
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

func (s SetlistService) HelperCreatePermSetlistImage(ctx context.Context, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- HelperCreatePermSetlistImage")

	if imageID == "" {
		log.Println("   No imageID provided")
		return "", errors.New("imageID empty")
	}

	sizes := map[string][2]int{
		"small":  {256, 192},
		"medium": {512, 384},
	}

	for name, _ := range sizes {

		sourceKey := fmt.Sprintf(
			"temp-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			name,
		)

		fmt.Println("sourcekey: ", sourceKey)

		destinationKey := fmt.Sprintf(
			"setlist-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
			imageID,
			name,
		)

		err := s.Storage.Copy(
			ctx,
			sourceKey,
			destinationKey,
		)
		if err != nil {
			log.Printf("   Unable to copy %s song image to R2 permanent storage: %v\n", name, err)
			return "", err
		}

		err = s.DeleteTempImage(ctx, imageID, bandSlug)
		if err != nil {
			log.Printf("   Unable to delete %s temporary setlist art: %v\n", name, err)
		}

	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"setlist-images",
		bandSlug,
		setlistSlug,
		imageID,
	)
	if err != nil {
		log.Println("   Unable to create browser path for setlist images: ", err)
		return "", err
	}

	return browserPath, nil
}

func (s SetlistService) DeleteTempImage(ctx context.Context, imageID string, bandSlug string) error {
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

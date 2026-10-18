package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/google/uuid"
)

func (s Service) SongsPage(ctx context.Context, user models.User, band models.Band) (models.MenuPageData, error) {

	songs, err := database.SongsTableGetAllSongsByBandID(ctx, band.BandID)
	if err != nil {
		return models.MenuPageData{}, fmt.Errorf("unable to get songs from DB: %w", err)
	}

	setlists, err := database.SetlistsTableGetSetlistsByBandIDAndUserID(ctx, band.BandID, user.UserID)
	if err != nil {
		return models.MenuPageData{}, fmt.Errorf("unable to get setlists: %w", err)
	}

	data := models.MenuPageData{
		User:     user,
		Band:     band,
		Songs:    songs,
		Setlists: setlists,
	}
	return data, nil
}

func (s Service) SongData(ctx context.Context, user models.User, band models.Band, songID string) (models.SongPageData, error) {

	song, err := database.SongsTableGetSongBySongID(songID, band.BandID)
	if err != nil {
		return models.SongPageData{}, fmt.Errorf("unable to get song from db: %w", err)
	}

	data := models.SongPageData{
		User: user,
		Band: band,
		Song: song,
	}

	return data, nil
}

func (s Service) SongDelete(ctx context.Context, user models.User, band models.Band, songID string) error {
	imageID, artworkPath, err := database.SongsTableGetImageIDAndPathBySongID(
		songID,
		band.BandID,
		user.UserID,
	)
	if err != nil {
		return fmt.Errorf("load authorized song artwork: %w", err)
	}

	if imageID != "" {
		if _, err := uuid.Parse(imageID); err != nil {
			return fmt.Errorf("invalid stored song artwork ID: %w", err)
		}

		localArtworkPath := "/static/uploads/song-images/" + imageID
		storedArtworkPath, pathErr := url.JoinPath(
			s.Storage.PublicURL,
			"song-images",
			band.Slug,
			imageID,
		)
		if pathErr != nil {
			return fmt.Errorf("build song artwork path: %w", pathErr)
		}

		switch artworkPath {
		case localArtworkPath:
			if err := s.ServiceDeleteITunesArtworkImageVersions(imageID); err != nil {
				return fmt.Errorf("delete imported song artwork: %w", err)
			}
		case storedArtworkPath:
			if err := s.ServiceDeleteArtworkImageVersions(ctx, imageID, band.Slug); err != nil {
				return fmt.Errorf("delete song artwork: %w", err)
			}
		default:
			log.Printf("Skipping unrecognized artwork path for song %s", songID)
		}
	}

	err = database.SongsTableDeleteSongByID(songID, band.BandID, user.UserID)
	if err != nil {
		return fmt.Errorf("delete song %s: %w", songID, err)
	}

	return nil
}

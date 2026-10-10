package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
)

func (s Service) SongsPage(ctx context.Context, user models.User, band models.Band) (models.MenuPageData, error) {

	songs, err := database.SongsTableGetAllSongsByBandID(band.BandID)
	if err != nil {
		return models.MenuPageData{}, errors.New("Unable to get songs from DB")
	}

	setlists, err := database.SetlistsTableGetSetlistsByBandIDAndUserID(band.BandID, user.UserID)
	if err != nil {
		return models.MenuPageData{}, fmt.Errorf("Unable to get setlists: ", err)
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

	song, err := database.SongsTableGetSongBySongID(songID)
	if err != nil {
		return models.SongPageData{}, fmt.Errorf("Unable to get song from db: %w", err)
	}

	data := models.SongPageData{
		User: user,
		Band: band,
		Song: song,
	}

	return data, nil
}

func (s Service) SongDelete(ctx context.Context, band models.Band, imageID, songID string) error {

	err := s.ServiceDeleteArtworkImageVersions(ctx, imageID, band.Slug)
	if err != nil {
		return fmt.Errorf("Unable to delete song image versions: ")
	}

	err = database.SongsTableDeleteSongByID(songID)
	if err != nil {
		return fmt.Errorf("Unable to delete songID: %s, err: %w", songID, err)
	}

	return nil
}
